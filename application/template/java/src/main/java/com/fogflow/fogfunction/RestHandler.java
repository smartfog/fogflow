package com.fogflow.fogfunction;

import java.io.IOException;
import java.net.URI;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import javax.annotation.PostConstruct;

import org.springframework.http.ResponseEntity;
import org.springframework.util.Assert;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import com.fasterxml.jackson.annotation.JsonAnySetter;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonGenerationException;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

class Config {
    Map<String, String> details = new LinkedHashMap<>();
 
    @JsonAnySetter
    void setDetail(String key, String value) {
        details.put(key, value);
    } 
}

class StatusCode {
	public int code;
	public String reasonPhrase;
	public String details;
	
	public StatusCode() {

	}		
	
	public int getCode() {
		return code;
	}	
	public void setCode(int code) {
		this.code = code;
	}			

	public String getReasonPhrase() {
		return reasonPhrase;
	}	
	public void setReasonPhrase(String reasonPhrase) {
		this.reasonPhrase = reasonPhrase;
	}	
	
	public String getDetails() {
		return details;
	}	
	public void setDetailsn(String details) {
		this.details = details;
	}		
}

class EntityId {
	public String id;
	public String type;
	public boolean isPattern;
	
	public EntityId() {

	}		
	
	public String getId() {
		return id;
	}	
	public void setId(String id) {
		this.id = id;
	}			

	public String getType() {
		return type;
	}	
	public void setType(String type) {
		this.type = type;
	}	
	
	public boolean getIsPattern() {
		return isPattern;
	}	
	public void setIsPattern(boolean isPattern) {
		this.isPattern = isPattern;
	}	
}

class ContextMetadata {
	public String name;
	public String type;
	public Object value;
	
	public ContextMetadata() {

	}		
	
	public String getName() {
		return name;
	}	
	public void setName(String name) {
		this.name = name;
	}			

	public String getType() {
		return type;
	}	
	public void setType(String type) {
		this.type = type;
	}	
	
	public Object getContextValue() {
		return value;
	}	
	public void setContextValue(Object value) {
		this.value = value;
	}		
}

class ContextAttribute {
	public String name;
	public String type;
	public Object contextValue;
	
	@JsonProperty("metadata")		
	public List<ContextMetadata> metadata;	
	
	public ContextAttribute() {

	}	
	
	public String getName() {
		return name;
	}	
	public void setName(String name) {
		this.name = name;
	}			

	public String getType() {
		return type;
	}	
	public void setType(String type) {
		this.type = type;
	}	
	
	public Object getContextValue() {
		return contextValue;
	}	
	public void setContextValue(Object contextValue) {
		this.contextValue = contextValue;
	}		
	
	@JsonIgnore
	public List<ContextMetadata> getMetadata() {
		if (metadata == null) {
			metadata = new ArrayList<ContextMetadata>();
		}
		
        return metadata;
    }

	@JsonIgnore
    public void setMetadata(List<ContextMetadata> metadata) {
        this.metadata = metadata;
    }  	
}


class ContextElement {
	public EntityId entityId;
	
	@JsonProperty("attributes")	
	public List<ContextAttribute> attributes;
	
	@JsonProperty("domainMetadata")	
	public List<ContextMetadata> domainMetadata;	
	
	public ContextElement() {
		attributes = new ArrayList<ContextAttribute>();
		domainMetadata = new ArrayList<ContextMetadata>();
	}	
	
	public ContextElement(ContextObject obj) {
		entityId = new EntityId();
		
		entityId.id = obj.id;
		entityId.type = obj.type;
		entityId.isPattern = false;
				
	    for (Map.Entry<String,ContextAttribute> entry : obj.attributes.entrySet()) {
			if (attributes == null) {
				attributes = new ArrayList<ContextAttribute>();
			}	    	
			attributes.add(entry.getValue());    	
	    }
	    
	    for (Map.Entry<String,ContextMetadata> entry : obj.domainMetadata.entrySet()) {
			if (domainMetadata == null) {
				domainMetadata = new ArrayList<ContextMetadata>();
			}	    	
			domainMetadata.add(entry.getValue());    	
	    }	    
	}		
	
	public EntityId getEntityId() {
		return entityId;
	}
	
	public void setEntityId(EntityId entityId) {
		this.entityId = entityId;
	}			
	
	@JsonIgnore	
	public List<ContextAttribute> getAttributes() {				
        return attributes;
    }

	@JsonIgnore
    public void setAttributes(List<ContextAttribute> attributes) {
        this.attributes = attributes;
    }	
    
	@JsonIgnore	
	public List<ContextMetadata> getDomainMetadata() {
        return domainMetadata;
    }
	
	@JsonIgnore
    public void setDomainMetadata(List<ContextMetadata> domainMetadata) {
        this.domainMetadata = domainMetadata;
    }    
}

class ContextElementResponse {
	public ContextElement  contextElement;
	public StatusCode  statusCode;
	
	public ContextElementResponse() {

	}		
	
	public ContextElement getContextElement() {
		return contextElement;
	}
	
	public void setContextElement(ContextElement contextElement) {
		this.contextElement = contextElement;
	}	
	
	public StatusCode getStatusCode() {
		return statusCode;
	}
	
	public void setStatusCode(StatusCode statusCode) {
		this.statusCode = statusCode;
	}		
}

class Notification {
	public String subscriptionId;
	public String originator;
	
	@JsonProperty("contextResponses")
	public List<ContextElementResponse> contextResponses;
	
	public Notification() {
		contextResponses = new ArrayList<ContextElementResponse>();
	}	
	
	public String getSubscriptionId() {
		return subscriptionId;
	}
	
	public void setSubscriptionId(String subscriptionId) {
		this.subscriptionId = subscriptionId;
	}	
	
	public String getOriginator() {
		return originator;
	}
	
	public void setOriginator(String originator) {
		this.originator = originator;
	}		
	
	@JsonIgnore
	public List<ContextElementResponse> getResponse() {	
        return contextResponses;
    }

	@JsonIgnore
    public void setResponse(List<ContextElementResponse> contextResponses) {
        this.contextResponses = contextResponses;
    }	
}


class ContextObject {
	public String id;
	public String type;
	
    public Map<String, ContextAttribute> attributes = new LinkedHashMap<>();
    
    public Map<String, ContextMetadata> domainMetadata = new LinkedHashMap<>();
    
	public ContextObject() {
		
	}

	public ContextObject(ContextElement element) {
		id = element.entityId.id;
		type = element.entityId.type;
		
		for(ContextAttribute attr : element.attributes) {
			attributes.put(attr.name,  attr);
		}
		
		for(ContextMetadata meta : element.domainMetadata) {
			domainMetadata.put(meta.name,  meta);
		}		
	}	    	
}


class UpdateContextRequest {
	public String updateAction;   
	
	@JsonProperty("contextElements")	
	public List<ContextElement> contextElements;
	
	public UpdateContextRequest() {
		contextElements = new ArrayList<ContextElement>();
	}	
		
	
	@JsonIgnore	
	public List<ContextElement> getContextElements() {				
        return contextElements;
    }

	@JsonIgnore
    public void setContextElements(List<ContextElement> contextElements) {
        this.contextElements = contextElements;
    }	
	
	public void addContextElement(ContextElement element) {
		this.contextElements.add(element);
	}
	
	public String getUpdateAction() {
		return updateAction;
	}
	
	public void setUpdateAction(String updateAction) {
		this.updateAction = updateAction;
	}	
}


class UpdateContextResponse {
	List<ContextElementResponse> contextResponses;  
	StatusCode   errorCode;
	
	public UpdateContextResponse() {
		contextResponses = new ArrayList<ContextElementResponse>();
	}	
		
	
	@JsonIgnore	
	public List<ContextElementResponse> getContextResponses() {				
        return contextResponses;
    }

	@JsonIgnore
    public void setContextResponses(List<ContextElementResponse> contextResponses) {
        this.contextResponses = contextResponses;
    }	
	
	
	public StatusCode getStatusCode() {
		return errorCode;
	}
	
	public void setStatusCode(StatusCode errorCode) {
		this.errorCode = errorCode;
	}	
}

@RestController
public class RestHandler {
	private String BrokerURL;
		
	private String outputEntityId;
	private String outputEntityType;
    
    @PostConstruct
    private void setup() {	
    	System.out.print("=========test========");
    	String jsonText = System.getenv("adminCfg");
    	this.handleInitAdminCfg(jsonText);    	
    }
	
	@PostMapping("/admin")
	public ResponseEntity<Void> handleConfig(@RequestBody List<Config> configs) {				
		for(Config cfg : configs) {
			System.out.println(cfg.details.get("command"));
			if (cfg.details.get("command").equalsIgnoreCase("CONNECT_BROKER") == true) {
				this.BrokerURL = cfg.details.get("brokerURL");				
			} else if (cfg.details.get("command").equalsIgnoreCase("SET_OUTPUTS") == true) {
				this.outputEntityId = cfg.details.get("id");
				this.outputEntityType = cfg.details.get("type");				
			}
		}
		
		System.out.println(this.BrokerURL.toString());
		System.out.println(this.outputEntityId.toString());		
		System.out.println(this.outputEntityType.toString());		
		
		return ResponseEntity.ok().build();
	}
	
	
	public void handleInitAdminCfg(String config)  {		
        ObjectMapper mapper = new ObjectMapper();
        
        try {
        	List<Config> myAdminCfgs =  mapper.readValue(config, new TypeReference<List<Config>>(){});
        	this.handleConfig(myAdminCfgs);        	
        } catch (JsonGenerationException e) {
			e.printStackTrace();
		} catch (JsonMappingException e) {
			e.printStackTrace();
		} catch (IOException e) {
			e.printStackTrace();
		}		
	}
	
	@PostMapping("/notifyContext")
	public ResponseEntity<Void> handleNotify(@RequestBody  Notification notify) {
		System.out.println(notify.subscriptionId);
				
		for(ContextElementResponse response : notify.contextResponses) {
			System.out.println(response.toString()); 
			if (response.statusCode.code == 200) {
				ContextObject contextObj = new ContextObject(response.contextElement);				
				handleEntityObject(contextObj);
			} 
		}	
		
		return ResponseEntity.ok().build();
	}	
	
	// to be overwritten according to your own data processing logic
	private void handleEntityObject(ContextObject entity) {
		System.out.println(entity.id);
		System.out.println(entity.type);
		
		ContextObject resultEntity = new ContextObject();
		resultEntity.id = "Test." + entity.id;
		resultEntity.type = "Result";
		
		ContextAttribute attr = new ContextAttribute();
		attr.name = "avg";
		attr.type = "integer";
		attr.contextValue = 10;
		
		resultEntity.attributes.put("avg", attr);
		
		publishResult(resultEntity);
	}
	

	// to publish the generated result entity		
	private void publishResult(ContextObject resultEntity) {		
		if (resultEntity.id == "") {
			resultEntity.id = outputEntityId;
		}
		
		if (resultEntity.type == "") {
			resultEntity.type = outputEntityType;
		}
		
		// send it to the nearby broker, assigned by FogFlow Broker
		RestTemplate restTemplate = new RestTemplate();
	    
		try {
		    URI uri = new URI(BrokerURL + "/updateContext");
		    
		    ContextElement element = new ContextElement(resultEntity);
		    
		    UpdateContextRequest request = new UpdateContextRequest();
		    request.addContextElement(element);
		    request.setUpdateAction("UPDATE");		    
		    
		    ObjectMapper mapper = new ObjectMapper();
	        try {
	            String json = mapper.writeValueAsString(request);
	            System.out.println("JSON = " + json);
	        } catch (JsonProcessingException e) {
	            e.printStackTrace();
	        }
		    	 
		    ResponseEntity<UpdateContextResponse> result = restTemplate.postForEntity(uri, request, UpdateContextResponse.class);
		    
		    System.out.println(result.getStatusCodeValue());		
		    
		}catch(Exception e) {
			e.printStackTrace();
		}
	}
}
