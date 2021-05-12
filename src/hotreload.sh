#!/bin/sh

ag -l | entr -r go run .

