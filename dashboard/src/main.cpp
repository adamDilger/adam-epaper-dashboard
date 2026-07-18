#define ENABLE_GxEPD2_GFX 0

#include <Arduino.h>

#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold9pt7b.h>
#include <time.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>

#include "graphics/GxEPD2_display_selection_new_style.h"
#include "wifi_creds.h"

static const char cloudflareRootCa[] PROGMEM = R"EOF(
-----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIQf+UwvzMTQ77dghYQST2KGzANBgkqhkiG9w0BAQsFADBX
MQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQMA4GA1UE
CxMHUm9vdCBDQTEbMBkGA1UEAxMSR2xvYmFsU2lnbiBSb290IENBMB4XDTIzMTEx
NTAzNDMyMVoXDTI4MDEyODAwMDA0MlowRzELMAkGA1UEBhMCVVMxIjAgBgNVBAoT
GUdvb2dsZSBUcnVzdCBTZXJ2aWNlcyBMTEMxFDASBgNVBAMTC0dUUyBSb290IFI0
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE83Rzp2iLYK5DuDXFgTB7S0md+8Fhzube
Rr1r1WEYNa5A3XP3iZEwWus87oV8okB2O6nGuEfYKueSkWpz6bFyOZ8pn6KY019e
WIZlD6GEZQbR3IvJx3PIjGov5cSr0R2Ko4H/MIH8MA4GA1UdDwEB/wQEAwIBhjAd
BgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zAd
BgNVHQ4EFgQUgEzW63T/STaj1dj8tT7FavCUHYwwHwYDVR0jBBgwFoAUYHtmGkUN
l8qJUC99BM00qP/8/UswNgYIKwYBBQUHAQEEKjAoMCYGCCsGAQUFBzAChhpodHRw
Oi8vaS5wa2kuZ29vZy9nc3IxLmNydDAtBgNVHR8EJjAkMCKgIKAehhxodHRwOi8v
Yy5wa2kuZ29vZy9yL2dzcjEuY3JsMBMGA1UdIAQMMAowCAYGZ4EMAQIBMA0GCSqG
SIb3DQEBCwUAA4IBAQAYQrsPBtYDh5bjP2OBDwmkoWhIDDkic574y04tfzHpn+cJ
odI2D4SseesQ6bDrarZ7C30ddLibZatoKiws3UL9xnELz4ct92vID24FfVbiI1hY
+SW6FoVHkNeWIP0GCbaM4C6uVdF5dTUsMVs/ZbzNnIdCp5Gxmx5ejvEau8otR/Cs
kGN+hr/W5GvT1tMBjgWKZ1i4//emhA1JG1BbPzoLJQvyEotc03lXjTaCzv8mEbep
8RqZ7a2CPsgRbuvTPBwcOMBBmuFeU88+FSBX6+7iP0il8b4Z0QFqIwwMHfs/L6K1
vepuoxtGzi4CZ68zJpiq1UvSqTbFJjtbD4seiMHl
-----END CERTIFICATE-----
)EOF";

#include <HTTPClient.h>

WiFiClientSecure client;
HTTPClient https;

void connectToWifi();
void screenMessage(const char *text);
bool syncClock();

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
  syncClock();
  client.setCACert(cloudflareRootCa);
}

int refreshCount = 0;
void loop()
{
  display.fillScreen(GxEPD_WHITE);

  if (refreshCount >= 10)
  {
    refreshCount = 0;
    fullRefresh();
  }
  else
  {
    refreshCount++;
  }

  int totalBytesRead = 0;

  if (https.begin(client, workerUrl))
  {
    https.addHeader("Accept", "application/octet-stream");

    int httpCode = https.GET();
    if (httpCode == 200)
    {
      int contentLength = https.getSize();
      if (contentLength <= 0)
      {
        Serial.printf("HTTP GET returned invalid content length: %d\n", contentLength);
        https.end();
        return;
      }

      Serial.printf("HTTP GET %d: %d bytes\n", httpCode, contentLength);

      uint8_t *responseBuffer = new uint8_t[contentLength];

      while (true)
      {
        if (totalBytesRead >= contentLength)
        {
          Serial.println("All content read.");
          break;
        }

        int bytesRead = https.getStream().readBytes(responseBuffer + totalBytesRead, contentLength - totalBytesRead);
        if (bytesRead <= 0)
        {
          Serial.println("HTTP stream ended before all content was read.");
          delete[] responseBuffer;
          https.end();
          return;
        }

        Serial.printf("%d total bytes read so far, %d bytes read this iteration\n", totalBytesRead, bytesRead);
        totalBytesRead += bytesRead;
      }

      for (int i = 0; i < totalBytesRead; i++)
      {
        byte b = responseBuffer[i];

        const int WIDTH = display.width();
        for (int j = 0; j < 8; j++)
        {
          int x = (i * 8 + j) % WIDTH;
          int y = (i * 8 + j) / WIDTH;

          if ((b & (1 << (7 - j))) != 0)
          {
            display.drawPixel(x, y, GxEPD_BLACK);
          }
        }
      }

      delete[] responseBuffer;

      https.end();
    }
    else
    {
      Serial.printf("HTTP GET failed, error: %d, message: %s\n", httpCode, https.errorToString(httpCode).c_str());
      https.end();
      return;
    }
  }

  partialRefresh();

  int tenMinutes = 10 * 60 * 1000;
  delay(tenMinutes);
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

bool syncClock()
{
  configTime(0, 0, "time.cloudflare.com", "pool.ntp.org");

  time_t now = time(nullptr);
  int attempts = 0;
  while (now < 1700000000 && attempts < 15)
  {
    delay(500);
    now = time(nullptr);
    attempts++;
  }

  if (now < 1700000000)
  {
    screenMessage("Clock sync failed.");
    return false;
  }

  return true;
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
