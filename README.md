Survana
=======

An HTML5 application for administering questionnaires on tablets and mobile devices. Developed by the Neuroinformatics Research Group at Harvard University.

Field:  button
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

        * option: @(option)

        * radiobutton:
            - id:string
            - name:string
            - html:string
            - value:string|number

        * checkboxbutton:



Html:
    - html:string

Size:
    - l: int = [1..12]
    - m: int = [1..12]
    - s: int = [1..12]
    - xs:int = [1..12]
