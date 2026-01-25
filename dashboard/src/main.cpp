#define ENABLE_GxEPD2_GFX 0

#include <Arduino.h>

#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold9pt7b.h>
#include <WiFi.h>

#include "graphics/GxEPD2_display_selection_new_style.h"
#include "network.h"
#include "wifi_creds.h"

WiFiClient client;

void connectToWifi();
void screenMessage(const char *text);

void testFunction()
{
  Serial.println("Test function called.");
  display.setTextColor(GxEPD_BLACK);
  display.setFont(&FreeMonoBold9pt7b);

  int HEIGHT = display.height();
  int WIDTH = display.width();

  int hello = 0;
  int halfHeight = HEIGHT / 2;
  int halfWidth = WIDTH / 2;

  while (true)
  {
    display.fillRect(0, 0, display.width(), display.height(), GxEPD_BLACK);
    Serial.printf("Starting loop iteration %d\n", hello);

    int squareSize = 40;
    // draw a checkerboard pattern with squares of size squareSize
    int rows = HEIGHT / squareSize;
    int cols = WIDTH / squareSize;
    for (int row = 0; row <= rows; row++)
    {
      for (int col = 0; col <= cols; col++)
      {
        if ((row + col) % 2 == 0)
        {
          int x = col * squareSize;
          int y = row * squareSize;
          display.fillRect(x, y, squareSize, squareSize, GxEPD_WHITE);
        }
      }
    }

    display.display(true);
    display.powerOff();
    hello++;
    delay(5000);
    Serial.println("Looping again.");
  }
}

void fullRefresh()
{
  display.fillScreen(GxEPD_WHITE);
  display.display(false);
  delay(1000);
}

void partialRefresh()
{
  display.display(true);
  display.powerOff();
}

void setup()
{
  display.init(115200, true, 2, false); // USE THIS for Waveshare boards with "clever" reset circuit, 2ms reset pulse
  // testFunction();
  connectToWifi();
}

int refreshCount = 0;
void loop()
{
  ResponseMetadata responseMetadata;

  display.fillScreen(GxEPD_WHITE);

  if (refreshCount >= 10)
    fullRefresh();
  else
    refreshCount++;

  doRequest(
      &responseMetadata,
      [](int x, int y, uint8_t count)
      {
        if (count == 1)
          display.drawPixel(x, y, GxEPD_BLACK);
        else
          display.drawLine(x, y, x + count - 1, y, GxEPD_BLACK);
      },
      client,
      display.width());

  partialRefresh();
  delay(responseMetadata.durationMinutes * 60 * 1000);
}

void connectToWifi()
{
  char msg[100];

  int retry = 0;

  sprintf(msg, "Attempting to connect to %s...", ssid);
  screenMessage(msg);

  WiFi.begin(ssid, pass);
  delay(1000);
  int status = WiFi.status();

  int attempts = 0;
  while (status != WL_CONNECTED)
  {
    attempts++;
    if (attempts > 5)
    {
      screenMessage("Failed to connect to wifi. Restarting...");
      ESP.restart();
    }

    sprintf(msg, "Couldn't get a wifi connection (error: %d). Retrying...", status);
    screenMessage(msg);
    delay(attempts * 1000);

    status = WiFi.status();
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
