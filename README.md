# Symulacja współbieżnej kraty w Go

Projekt powstał jako zabawa Go i jego prostym modelem programowania współbieżnego. Program uruchamia symulację na prostokątnej kracie, po której poruszają się ludzie oznaczeni numerami. W tle, niezależnie od siebie, pojawiają się też zagrożenia i dzicy lokatorzy, a całość jest jednocześnie renderowana w terminalu.

## Co robi program

- generuje planszę o rozmiarze `30 x 10`,
- tworzy ludzi oznaczonych identyfikatorami `00`, `01`, `02`, ...,
- porusza każdego człowieka losowo po sąsiednich polach,
- co jakiś czas tworzy zagrożenia `##`,
- co jakiś czas tworzy dzikich lokatorów `**`,
- odświeża widok planszy niezależnie od logiki ruchu.

## Zasady symulacji

- Człowiek porusza się losowo na jedno z dostępnych pól sąsiednich.
- Jeśli człowiek wejdzie na pole z zagrożeniem `##`, znika z planszy.
- Dziki lokator `**` próbuje uciec, gdy ktoś chce wejść na jego pole.
- Jeśli lokator nie ma gdzie uciec, zostaje usunięty i jego miejsce zajmuje nadchodząca postać.
- Zagrożenia i lokatorzy żyją tylko przez ograniczony czas, po czym same znikają.
- Krawędzie planszy są zablokowane, więc ruch odbywa się tylko wewnątrz kraty.

## Oznaczenia na planszy

- `00`, `01`, `02`, ...: ludzie
- `##`: zagrożenie
- `**`: dziki lokator
- `-` i `|`: ślady ostatnich ruchów widoczne pomiędzy polami
- puste pole: brak obiektu

## Współbieżność

Najważniejsza idea projektu to rozdzielenie odpowiedzialności między wiele goroutine i komunikację przez kanały:

- każde pole kraty działa jak niezależny węzeł obsługujący własny stan,
- ruch między sąsiadami jest realizowany komunikatami wysyłanymi kanałami,
- osobny `Reviver` odpowiada za pojawianie się nowych ludzi, lokatorów i zagrożeń,
- osobny `Snapshotter` zbiera aktualizacje z planszy i rysuje stan świata w terminalu,
- logika symulacji i wyświetlanie działają równolegle.

## Parametry czasowe

Domyślna konfiguracja jest zapisana w pliku [server_solution/consts.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/consts.go):

- odświeżanie planszy: co `2s`,
- nowy człowiek: co `6s`,
- nowy lokator: co `5s`,
- nowe zagrożenie: co `10s`,
- czas życia lokatora: `12s`,
- czas życia zagrożenia: `17s`,
- ruch człowieka: losowo co `3-4s`.

## Uruchomienie

Wymagany jest Go `1.21.2` lub nowszy.

```bash
go run .
```

Program działa do momentu naciśnięcia `ENTER`.

## Struktura projektu

- [main.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/main.go): start programu
- [server_solution/init.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/init.go): budowa kraty i kanałów
- [server_solution/vertex.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/vertex.go): logika pojedynczego pola, ruch, kolizje
- [server_solution/reviver.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/reviver.go): tworzenie nowych bytów
- [server_solution/snapshotter.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/snapshotter.go): render planszy

## Co warto pokazać na demo

- start z pustą planszą,
- pojawienie się pierwszego człowieka i jego ścieżki ruchu,
- pojawienie się lokatora `**` i jego ucieczkę,
- wejście człowieka na `##` i usunięcie go z planszy,
- równoległe logi tekstowe i odświeżanie widoku planszy.

## Uwagi

- Przy `AreCommentLabels = true` program wypisuje dodatkowe komunikaty diagnostyczne w terminalu.
- Parametry symulacji można łatwo zmieniać przez stałe w [server_solution/consts.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/consts.go).
