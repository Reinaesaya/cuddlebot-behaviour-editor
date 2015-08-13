# Cuddlebot Behaviour Editor

Adapted from [JeffBehaviour Editor](https://github.com/lauracang/JeffBehaviourEditor-master)

This updated version includes various command line functions written in Javascript the facilitate smooth actuator motion, as well as applied, complex head motions. It does not improve the UI of the behaviour editor.

## Setup

The following set of instructions applies for Windows 7, and may or may not be applicable to other operating systems

Step 1: Install [go](https://golang.org/doc/install) and [git](https://git-scm.com/downloads)
Step 2: Clone folder via git
```bash
git clone https://github.com/reinaesaya/cuddlebot-behaviour-editor
```
Step 3: Navigate to folder via bash or cmd, and run editor
```bash
go run behaviourServer.go
```
Step 4: Open up a browser to `localhost:8080`
Step 5: Right click on a blank space on the page and click `Inspect element`. Switch to `console`

## Usage

### Visual

Basic GUI features can be tested out via trial and error. It should be noted however, that setting the loop to 65535 leads to the action be done an infinite amount of times. This can be stopped through a command line sleep instruction for the particular actuator, demonstrated as follows:
```javascript
sendSleepCommand(["headx"]); // sleep left-right head direction
```

### Command Line

There are a variety of command line functions that are not presented in the GUI. These commands can be run, and used either as baseline functions in further development or altered and referenced.

The following are some examples of the command line functions able to be used and implemented.

```javascript
// Move head to left-most point and stay there for 1 second
sendSetPointCommand("headx", 0, 1, [1000,0]);

// Move head up for 1 second then down for 1 second, 3 times
sendSetPointCommand("heady", 0, 3, [1000,65535,1000,0]);

// Move head from current point to left-most point in 2 seconds,
// 	and stay there for 3 seconds
sendSmoothCommand("headx", 2000, [3000,0]);

// Purr infinitely long
sendSetPointCommand("purr", 0, 65535, [1231,35421]);

// Stop purring
sendSleepCommand(["purr"]); // Note the brackets


// Basic breathing (1 second out, 1 second relaxed)
startBreathing(2000);

// Stop breathing
stopBreathing();


// Move head from left-most point to right-most point in 2 seconds
smoothToPt("headx", 0, 65535, 2000);

// Rotate head from top right corner to 2 rotations counterclockwise in 5 seconds
rotateHead(30000, 0, 720, "Counterclockwise", 5000)
```
