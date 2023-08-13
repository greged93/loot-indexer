# Loot indexer

This repository is a work in progress for an indexer for the loot survivor game.
The goal is to index and display in an interactive way the game progress of the
top 3 to 5 players.

# TODO

-   [ ] Update the way events are stored:
    -   only store the raw bytes of the event.
    -   store the evolution of the adventurer based on reserved events (e.g.
        DiscoveredGold -> increase Adventurer's gold).
