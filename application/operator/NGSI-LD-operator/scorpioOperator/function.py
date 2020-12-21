ctxELement = {}

def handleEntity(ctxObj, create, update, append):
    print('===============Implement losic====================')
    print(ctxObj)
    for ctx in ctxObj:
	handleScorpioUpdate(ctx, create, update, append)


def handleupdateAppend(currUpdateCtx, create, update, append):    
    appendCtx = {}
    global ctxELement
    eid = currUpdateCtx['id']
    preCtxEle = ctxELement[eid]
    for key in  currUpdateCtx : 
        if ctxELement[eid].has_key(key) == False :
	    appendCtx[key] = currUpdateCtx[key]
	    preCtxEle[key] = currUpdateCtx[key]
    
    ctxELement[eid] = preCtxEle
    print("This is ctxELement")
    if len(appendCtx) > 0: 
        append(appendCtx)
    update(currUpdateCtx)
    

def handleScorpioUpdate(ctx, create, update, append):
	global ctxELement 
	eid = ctx['id']
	print(eid)
	if ctxELement.has_key(eid) == True:
	    handleupdateAppend(ctx, create, update, append)
	else:
	    ctxELement[eid] = ctx
	    create(ctx)
    
    
