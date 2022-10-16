# Password-Gen
Met behulp van deze applicatie kun je een random wachtwoord laten generaten op basis van je eigen wensen, dit wachtwoord wordt vervolgens in een database gezet.

## Benodigdheden
Deze applicatie werkt met de versie 1.19 van go. Verder worden er 3 verschillende packages gebruikt die gedownload moeten worden.
- github.com/lib/pq v1.10.7
- github.com/sethvargo/go-password v0.2.0 
- gopkg.in/yaml.v3 v3.0.1

## Config File
In de applicatie wordt er een yaml file gebruikt voor configuratie van de database, hierin moet de inlog van de database staan zodat de applicatie hierbij kan. Zorg dat deze juist is ingevuld bij gebruik van applicatie.

## Flags
Deze applicatie maakt gebruikt van meerdere flags, deze flags kun je aanpassen in de commandline bij het gebruik van de applicatie.
Int Length "-l":    De eerste flag bepaalt de lengte van het wachtwoord wat gegenereerd wordt.
                    Deze moet worden ingevuld anders zal de applicatie niet werken.
Bool Digits "-digits":  Dit is een flag waarbij je kunt bepalen of er nummers bij het wachtwoord
                        gebruikt mogen worden. Dit staat automatisch op true wil je dit niet zet hem naar false. 
Bool Symbols "-symbols":    Deze flag wordt gebruikt om te bepalen of er symbolen in je wachtwoord
                            mogen zitten. Deze staat automatisch op true, zet deze naar false wanneer je dit niet wilt.
Bool Lower "-lower":    Deze bool wordt gebruikt om aan te geven of er kleine letters in het
                        wachtwoord mogen zitten. Deze staat op true, zet naar false als je dit niet wilt
Bool Repeat "-repeat":  Deze bool wordt gebruikt om aan te geven of je wilt dat characters zich mogen
                        herhalen in het wachtwoord. Deze staat op true, zet naar false als je dit niet wilt.

## Voorbeeld
Hieronder staat een voorbeeld van hoe je de applicatie kunt gebruiken.
```
go run. -l 8
```
Dit maakt een wachtwoord van 8 characters en stuurt deze naar de database aangegeven in de config file.
