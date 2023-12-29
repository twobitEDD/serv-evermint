### How to make your own EVM-enabled chain from this fork:
_(support Linux & MacOS only)_

Step 1: Open `constants.go` and change every constant in there.

Step 2: Do `git commit` (+ `git push`) to save your changes.

Step 3: Run `./rename-chain.sh`

Done, you have your own chain now.
___
Cleanup notes to be checked after running script:
- Make sure directory `./rename_chain` is deleted.
- Make sure `./rename-chain.sh` is deleted.
