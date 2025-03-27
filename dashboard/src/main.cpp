#define ENABLE_GxEPD2_GFX 0

#include <Arduino.h>

#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold9pt7b.h>

#include "graphics/GxEPD2_display_selection_new_style.h"

void helloWorld();

void setup()
{
  Serial.begin(115200);
  Serial.println();
  Serial.println("setup");
  delay(100);
  display.init(115200, true, 2, false); // USE THIS for Waveshare boards with "clever" reset circuit, 2ms reset pulse
  helloWorld();
}

void loop()
{
  // put your main code here, to run repeatedly:
}

const char HelloWorld[] = "Hello World!";
const char HelloArduino[] = "Hello Arduino!";
const char HelloEpaper[] = "Hello E-Paper!";

void helloWorld()
{
  // Serial.println("helloWorld");
  display.setRotation(0);
  display.setFont(&FreeMonoBold9pt7b);
  if (display.epd2.WIDTH < 104)
    display.setFont(0);
  display.setTextColor(GxEPD_BLACK);
  int16_t tbx, tby;
  uint16_t tbw, tbh;
  display.getTextBounds(HelloWorld, 0, 0, &tbx, &tby, &tbw, &tbh);
  // center bounding box by transposition of origin:
  uint16_t x = ((display.width() - tbw) / 2) - tbx;
  uint16_t y = ((display.height() - tbh) / 2) - tby;
  display.setFullWindow();
  display.firstPage();
  do
  {
    display.fillScreen(GxEPD_WHITE);
    display.setCursor(x, y);
    display.print(HelloWorld);
  } while (display.nextPage());
  // Serial.println("helloWorld done");
}