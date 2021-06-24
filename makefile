run: jarmuzrgblight
	./jarmuzrgblight

jarmuzrgblight: *.go
	go build -o jarmuzrgblight .

build:
	go build -o jarmuzrgblight .
