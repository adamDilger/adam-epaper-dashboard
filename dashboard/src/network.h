#include <WiFi.h>
#include <stdint.h>

const char url[] = "dashboard.dlgr.au";
const char endOfHeaders[] = "\r\n\r\n";

struct ResponseMetadata
{
    uint8_t formatVersion;
    uint8_t durationMinutes;
};

int doRequest(
    ResponseMetadata *responseMetadata,
    std::function<void(int, int, uint8_t)> drawLine,
    WiFiClient client,
    int screenWidth);
