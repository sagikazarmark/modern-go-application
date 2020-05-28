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


    Scenario: Adding a new item without a title
        Given an empty todo list
        When the user adds a new item for ""
        Then it should fail with a validation error for the "title" field saying that "title cannot be empty"
        And the list should be empty


    Scenario: An item can be marked as complete
        Given an empty todo list
        And an item for "Call mom"
        When it is marked as complete
        Then it should be complete
        But it should be on the list
