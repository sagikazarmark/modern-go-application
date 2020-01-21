Feature: Todo list

    Scenario: Add a new item to the list
        Given there is an empty todo list
        When I add entry "Call mom"
        Then I should have a todo to "Call mom"

    Scenario: Cannot add an empty todo item
        Given there is an empty todo list
        When I add entry ""
        Then I should see a validation error for the "text" field saying that "text cannot be empty"
        #And I should not have a todo to "Call mom"
#
#    Scenario: Mark a todo item as done
#        Given there is an empty todo list
#        And the entry "Call mom" is on the list
#        When I mark "Call mom" as done
#        Then "Call mom" entry should appear as done
