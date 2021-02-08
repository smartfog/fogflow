const config = {};

function toBoolean(env, defaultValue) {
  return env !== undefined ? env.toLowerCase() === 'true' : defaultValue;
}

function to_array(env, default_value){
    return (env !== undefined) ? env.split(',') : default_value;
}

// Used only if https is disabled
config.pep_port = process.env.PEP_PROXY_PORT || 5555;

// Set this var to undefined if you don't want the server to listen on HTTPS
config.https = {
  enabled: toBoolean(process.env.PEP_PROXY_HTTPS_ENABLED, false),
  cert_file: 'cert/cert.crt',
  key_file: 'cert/key.key',
  port: process.env.PEP_PROXY_HTTPS_PORT || 443,
};

//IDM configration
config.idm = {
  host: process.env.PEP_PROXY_IDM_HOST || '',
  port: process.env.PEP_PROXY_IDM_PORT || 3000,
  ssl: toBoolean(process.env.PEP_PROXY_IDM_SSL_ENABLED, false),
};

//Apllication information, for which account has been registered on IDM
config.app = {
  host: process.env.PEP_PROXY_APP_HOST || '',
  port: process.env.PEP_PROXY_APP_PORT || '80',
  ssl: toBoolean(process.env.PEP_PROXY_APP_SSL_ENABLED, false), // Use true if the app server listens in https
};

config.organizations = {
  enabled: toBoolean(process.env.PEP_PROXY_ORG_ENABLED, false),
  header: process.env.PEP_PROXY_ORG_HEADER || 'fiware-service'
}

// Credentials obtained when registering PEP Proxy in app_id in Account Portal
config.pep = {
  app_id: process.env.PEP_PROXY_APP_ID || '',
  username: process.env.PEP_PROXY_USERNAME || '',
  password: process.env.PEP_PASSWORD || '',
  token: {
    secret: process.env.PEP_TOKEN_SECRET || '', // Secret must be configured in order validate a jwt
  },
  trusted_apps: [],
};

// in seconds
config.cache_time = 31557600;

// if enabled PEP checks permissions in two ways:
//  - With IdM: only allow basic authorization
//  - With Authzforce: allow basic and advanced authorization.
//        For advanced authorization, you can use custom policy checks by including programatic scripts
//    in policies folder. An script template is included there
//
//      This is only compatible with oauth2 tokens engine

config.authorization = {
  enabled: toBoolean(process.env.PEP_PROXY_AUTH_ENABLED, false),
  pdp: process.env.PEP_PROXY_PDP || 'idm', // idm|authzforce
  azf: {
    protocol: process.env.PEP_PROXY_AZF_PROTOCOL || 'http',
    host: process.env.PEP_PROXY_AZF_HOST || 'localhost',
    port: process.env.PEP_PROXY_AZF_PORT || 8080,
    custom_policy: process.env.PEP_PROXY_AZF_CUSTOM_POLICY || undefined, // use undefined to default policy checks (HTTP verb + path).
  },
};

// list of paths that will not check authentication/authorization
// example: ['/public/*', '/static/css/']
config.public_paths = to_array(process.env.PEP_PROXY_PUBLIC_PATHS, []);

config.magic_key = process.env.PEP_PROXY_MAGIC_KEY || undefined;

module.exports = config;

