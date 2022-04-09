from threading import *
ctxELement = {}


def handleEntity(ctxObj, create, update, append):
    print('===============Implement losic====================')
    handleScorpioUpdate(ctxObj, create, update, append)


# handle update and append request for scorpio broker

def handleupdateAppend(currUpdateCtx, create, update, append):
    appendCtx = {}
    global ctxELement
    eid = currUpdateCtx['id']
    preCtxEle = ctxELement[eid]
    appendCtx['id'] = currUpdateCtx['id']
    appendCtx['type'] = currUpdateCtx['type']
    for key in currUpdateCtx:
        if ctxELement[eid].has_key(key) == False:
            appendCtx[key] = currUpdateCtx[key]
            preCtxEle[key] = currUpdateCtx[key]

    ctxELement[eid] = preCtxEle
    if len(appendCtx) > 2:
        ctxELement[eid] = preCtxEle
        append(appendCtx)
    update(currUpdateCtx)


# handle creation of etity on scorpio broker

def handleScorpioUpdate(ctx, create, update, append):
    global ctxELement
    eid = ctx['id']
    if ctxELement.has_key(eid) == True:
        handleupdateAppend(ctx, create, update, append)
    else:
        ctxELement[eid] = ctx
        create(ctx)

# handle  case if entity is already persent on scorpio broker


def handleAlreadyCreatedEntity(eid, create, update, append):
    global ctxELement
    ctxObj = ctxELement[eid]
    handleupdateAppend(ctxObj, create, update, append)
