#include <network.h>

const int BUFFER_SIZE = 128;
uint8_t responseBuffer[BUFFER_SIZE] = {0};

int doRequest(
    ResponseMetadata *responseMetadata,
    std::function<void(int, int, uint8_t)> drawLine,
    WiFiClient client,
    int screenWidth)
{
    if (!client.connect(url, 80))
    {
        Serial.println("Connection failed");
        delay(10000);
        return 1;
    }

    // Send HTTP request
    client.println("GET / HTTP/1.1");
    client.printf("Host: %s\n", url);
    client.println("Accept: application/octet-stream");

    if (client.println() == 0)
    {
        Serial.println("Failed to send request");
        client.stop();
        delay(10000);
        return 1;
    }

    // Check HTTP status
    char status[32] = {0};
    client.readBytesUntil('\r', status, sizeof(status));
    if (strcmp(status, "HTTP/1.1 200 OK") != 0)
    {
        Serial.print("Unexpected response: ");
        Serial.printf("[%s]", status);
        Serial.println();
        client.stop();
        delay(10000);
        return 1;
    }

    if (!client.find("Content-Length: "))
    {
        Serial.println("Invalid response, no content length");
        client.stop();
        delay(10000);
        return 0;
    }

    int contentLength = client.parseInt();
    if (contentLength <= 0 || contentLength > 48000)
    {
        Serial.println("Invalid response, content length is 0 or too large");
        client.stop();
        delay(10000);
        return 0;
    }

    // Skip HTTP headers
    if (!client.find(endOfHeaders))
    {
        Serial.println("Invalid response");
        client.stop();
        delay(10000);
        return 0;
    }

    Serial.printf("drawing bitmap with content length of %d\n", contentLength);

    client.read(&responseMetadata->formatVersion, 1);
    client.read(&responseMetadata->durationMinutes, 1);

    Serial.printf("format version: %d\n", responseMetadata->formatVersion);
    Serial.printf("duration minutes: %d\n", responseMetadata->durationMinutes);

    int y = 0;
    int x = 0;
    uint8_t count = 0;
    bool isBlack = false;

    int bufferSize = BUFFER_SIZE;

    while (client.available() > 0)
    {
        while (true)
        {
            int av = client.available();
            if (av == 0)
            {
                break;
            }

            if (av < bufferSize)
            {
                bufferSize = av;
            }

            client.read(responseBuffer, bufferSize);

            for (int pos = 0; pos < bufferSize; pos++)
            {
                isBlack = (responseBuffer[pos] & 0b10000000) != 0;
                count = responseBuffer[pos] & 0b01111111;

                if (isBlack)
                {
                    // display.drawLine(x, y, x + count, y, GxEPD_BLACK);
                    drawLine(x, y, count);
                }

                x += count;

                if (x >= screenWidth)
                {
                    x = 0;
                    y++;
                }
            }
        }

        // TODO: this is a bit rubbish, wait 800ms between each read chunk
        // to allow for more data to be available, then start reading again
        // if we are still waiting for data
        delay(800);
    }

    Serial.printf("fin\n");

    client.stop();

    return 0;
}