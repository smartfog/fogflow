function registerAllBlocks(blocks){

blocks.register({
    name: "FogFunction",
    description: "To specify a fog function",
    fields: [
        {
            name: "Name",
            type: "string",
            defaultValue: "unknown",
            attrs: "editable"
        },     
        {
            name: "User",
            type: "string",
            defaultValue: "fogflow",
            attrs: "editable"
        },    
        {
            name: "Selectors",
            type: "Selector",
            attrs: "input"
        },
        {
            name: "Annotators",
            type: "Annotator",
            attrs: "output"
        }
    ]
});

blocks.register({
    name: "InputTrigger",
    description: "To select a type of input stream of a fog function",
    fields: [
        {
            name: "SelectedAttributes",
            type: "string[]",
            defaultValue: ["all"],
            attrs: "editable"
        },          
        {
            name: "Groupby",
            type: "string[]",
            defaultValue: ["all"],            
            attrs: "editable"
        },              
        {
            name: "Conditions",
            type: "Condition",
            attrs: "input"
        },        
        {
            name: "Selector",
            type: "Selector",
            attrs: "output"
        }        
    ]
});

blocks.register({
    name: "SelectCondition",
    description: "To define a restriction for entity selection",
    fields: [
        {
            name: "Type",
            type: "choice",
            choices: ["EntityId","EntityType", "GeoScope(Nearby)", "GeoScope(InCircle)", "GeoScope(InPolygon)", "TimeScope", "StringQuery"],
            defaultValue: "EntityType",
            attrs: "editable"
        },
        {
            name: "value",
            type: "string",
            attrs: "editable"
        },   
        {
            name: "Condition",
            type: "Condition",
            attrs: "output"
        }             
    ]
});

blocks.register({
    name: "OutputAnnotator",
    description: "To annotate the output stream of a fog function",
    fields: [
        {
            name: "EntityType",
            type: "string",
            attrs: "editable"
        },
        {
            name: "Herited",
            type: "bool",
            defaultValue: false,
            attrs: "editable"
        },
        {
            name: "Annotator",
            type: "Annotator",
            attrs: "input"
        }        
    ]
});

}



