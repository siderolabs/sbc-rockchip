# sbc-template

This is a template repository to generate a new project for the SBC.

The template already contains an example of building and copying over `u-boot`.
This can be used as a starting point for your own project.

## Requirements

Using a sbc project templates assume docker is installed on your machine.

## Usage

* Create a new repository using this template.
* Start by replacing `board` to `<your-board-name>`
* Update `installers/<your-board-name>/src/main.go` with your own logic.
* Run `make rekres`
* Now run `make help` on instructions on setting up `buildx` builder if one is not already setup.
