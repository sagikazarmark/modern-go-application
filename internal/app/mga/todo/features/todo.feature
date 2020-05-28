Feature: Todo list


    Scenario: Add a new todo to the list
        When the user adds a new todo "Call mom" to the list
        Then "Call mom" should be on the list


    Scenario: Cannot add an empty todo
        When the user adds a new todo "" to the list
        Then it should fail with a validation error for the "title" field saying that "title cannot be empty"


    Scenario: Mark a todo as complete
        Given there is a todo "Call mom"
        When the user marks it as complete
        Then it should be complete
