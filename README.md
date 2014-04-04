Survana
=======

An HTML5 application for administering questionnaires on tablets and mobile devices. Developed by the Neuroinformatics Research Group at Harvard University.

Field:
------
        button
        input
        text
        label
        instructions
        separator
        select
        group
        matrix

    * button:
        - html:string
        - click:string //JS function name
        - align:string = ["center", "left", "right"]

    * input:
        - id:string
        - name:string
        - placeholder:string
        - label: <Field:label>
        - size: <Size>
        - prefix: <Field:*>
        - suffix: <Field:*>

    * text: @(input)

    * label:
        - html:string
        - align:string = ["center", "left", "right"]
        - size: <Size>

    * instructions:
        - html:string
        - align:string = ["center", "left", "right"]

    * separator

    * option:
        - html:string
        - value:string|number

    * select
        - id:string
        - name:string
        - fields: Array(<Field:option,group:option>)

    * group
        - id:string
        - name:string
        - label: <Field:label>
        - group:string = <GroupField:*>
        - fields: Array(<GroupField:*>)

    * matrix
        - striped:bool
        - hover:bool
        - equalize:bool
        - matrix:string = <Field:*> = Field:radio
        - columns:Array(<Column>)
        - rows:Array(<Row>)


GroupField:
----------
            button
            radio
            checkbox
            option

        * button:
            - id:string
            - name:string
            - html:string
            - value:string|number

        * radio:
            - id:string
            - name:string
            - html:string
            - value:string|number

        * checkbox:
            - id:string
            - name:string
            - html:string
            - value:string|number

        * option:       @(option)

        * radiobutton:
            - id:string
            - name:string
            - html:string
            - value:string|number

        * checkboxbutton:
            - id:string
            - name:string
            - html:string
            - value:string|number


Html:
    - html:string

Size:
    - l: int = [1..12]
    - m: int = [1..12]
    - s: int = [1..12]
    - xs:int = [1..12]
