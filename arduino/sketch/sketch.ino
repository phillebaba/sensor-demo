#include <OneWire.h>
#include <DallasTemperature.h>

#define ONE_WIRE_BUS 2

#define SOP '<'
#define EOP '>'

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

void setup(void)
{
  Serial.begin(9600);
  sensors.begin();
}


void loop(void)
{
  sensors.requestTemperatures(); 
  float temperature = sensors.getTempCByIndex(0);
  String temperatureString = String(temperature);
  
  Serial.print(SOP);
  Serial.print(temperatureString);
  Serial.print(EOP);
  delay(200);
}
