Feature: Todo list

# Inspired by https://paulhammant.com/2017/05/14/todomvc-and-given-when-then-scenarios/
# and http://todobackend.com/


    Scenario: Adding a new item to an empty list
        Given an empty todo list
        When the user adds a new item for "Call mom"
        Then it should be the only item on the list


    Scenario: Adding new items to an empty list
        Given an empty todo list
        When the user adds a new item for "File taxes"
        And the user also adds a new item for "Walk the dog"
        Then both items should be on the list


    Scenario: An item can be marked as complete
        Given an empty todo list
        And an item for "Call grandma"
        When it is marked as complete
        Then it should be complete
        But it should be on the list


    Scenario: An item can be deleted
        Given an empty todo list
        And an item for "Buy milk"
        When it is deleted
        Then the list should be empty


    Scenario: All items can be deleted
        Given an empty todo list
        And an item for "Buy milk"
        And an item for "Buy cheese"
        When all items are deleted
        Then the list should be empty
