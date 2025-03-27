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

const int RESPONSE_SIZE = 48000;

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

  doRequest(url);

  display.display(false);
}

void loop()
{
  // put your main code here, to run repeatedly:
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

  // Skip HTTP headers
  if (!client.find(endOfHeaders))
  {
    Serial.println("Invalid response");
    client.stop();
    delay(10000);
    return 0;
  }

  Serial.printf("drawing bitmap\n");

  int w_byte_count = display.width() / 8;
  uint8_t *response = (uint8_t *)malloc(w_byte_count);

  // read each row of the response one by one, and draw it on the screen
  for (int y = 0; y < display.height(); y++)
  {
    client.read(response, w_byte_count);

    // for each byte in the row
    for (int pos = 0; pos < w_byte_count; pos++)
    {
      // loop over each bit, drawing a pixel if it's set
      for (int i = 0; i < 8; i++)
      {
        if ((response[pos] & (1 << i)) != 0)
        {
          display.drawPixel((pos * 8) + (7 - i), y, GxEPD_BLACK);
        }
      }
    }
  }

  Serial.printf("fin\n");

  client.stop();

  return 0;
}
