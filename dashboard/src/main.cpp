#define ENABLE_GxEPD2_GFX 0

#include <Arduino.h>

#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold9pt7b.h>
#include <WiFi.h>

#include "graphics/GxEPD2_display_selection_new_style.h"
#include "wifi_creds.h"

WiFiClient client;
int status = WL_IDLE_STATUS;
char url[] = "192.168.0.147";
const char endOfHeaders[] = "\r\n\r\n";

void connectToWifi();
void screenMessage(const char *text);
int doRequest(const char *url);

void setup()
{
  Serial.begin(115200);
  display.init(115200, true, 2, false); // USE THIS for Waveshare boards with "clever" reset circuit, 2ms reset pulse
  screenMessage("Hello World!");
  connectToWifi();

  display.fillScreen(GxEPD_WHITE);
  display.display(true);
}

void loop()
{
  while (true)
  {
    display.fillScreen(GxEPD_WHITE);
    doRequest(url);
    display.display(true);
    delay(8000);
  }
}

void connectToWifi()
{
  char msg[100];

  int retry = 0;

  sprintf(msg, "Attempting to connect to %s...", ssid);
  screenMessage(msg);

  WiFi.begin(ssid, pass);
  delay(1000);
  status = WiFi.status();

  if (status != WL_CONNECTED)
  {
    sprintf(msg, "Couldn't get a wifi connection (error: %d)\nPlease reboot.", status);
    screenMessage(msg);
    while (true)
      ;
  }

  screenMessage("Connected to wifi.");
}

void screenMessage(const char *text)
{
  Serial.printf("Drawing message: [%s]\n", text);
  display.setRotation(0);
  display.setFont(&FreeMonoBold9pt7b);
  display.setTextColor(GxEPD_BLACK);

  int16_t tbx, tby;
  uint16_t tbw, tbh;
  display.getTextBounds(text, 0, 0, &tbx, &tby, &tbw, &tbh);

  // center bounding box by transposition of origin:
  uint16_t x = ((display.width() - tbw) / 2) - tbx;
  uint16_t y = ((display.height() - tbh) / 2) - tby;

  display.setFullWindow();

  display.fillScreen(GxEPD_WHITE);
  display.setCursor(x, y);
  display.print(text);
  display.display(true);
}

int doRequest(const char *url)
{
  if (!client.connect(url, 8000))
  {
    Serial.println("Connection failed");
    delay(10000);
    return 1;
  }

  // Send HTTP request
  client.println("GET / HTTP/1.1");
  client.println("Host: 192.168.0.147");

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

  int bufferSize = 128;
  uint8_t *response = (uint8_t *)malloc(bufferSize);

  int y = 0;
  int x = 0;
  uint8_t count = 0;
  bool isBlack = false;

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

    client.read(response, bufferSize);

    for (int pos = 0; pos < bufferSize; pos++)
    {
      isBlack = (response[pos] & 0b10000000) != 0;
      count = response[pos] & 0b01111111;

      if (isBlack)
      {
        display.drawLine(x, y, x + count, y, GxEPD_BLACK);
      }

      x += count;

      if (x >= display.width())
      {
        x = 0;
        y++;
      }
    }
  }

  Serial.printf("fin\n");

  client.stop();

  return 0;
}
