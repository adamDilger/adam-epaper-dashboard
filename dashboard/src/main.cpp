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

void setup()
{
  Serial.begin(115200);
  display.init(115200, true, 2, false); // USE THIS for Waveshare boards with "clever" reset circuit, 2ms reset pulse
  connectToWifi();

  display.fillScreen(GxEPD_WHITE);
  display.display(true);
}

int refreshCount = 0;

void loop()
{
  ResponseMetadata responseMetadata;

  display.fillScreen(GxEPD_WHITE);

  doRequest(
      &responseMetadata,
      [](int x, int y, uint8_t count)
      {
        display.drawLine(x, y, x + count, y, GxEPD_BLACK);
      },
      client,
      display.width());

  if (refreshCount > 10)
  {
    display.display(false);
    refreshCount = 0;
  }
  else
  {
    display.display(true);
    refreshCount++;
  }

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
