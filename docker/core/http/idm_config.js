const config = {};

function to_boolean(env, default_value){
    return (env !== undefined) ? (env.toLowerCase() === 'true') : default_value;
}

function to_array(env, default_value){
    return (env !== undefined) ? env.split(',') : default_value;
}

config.port = (process.env.IDM_PORT || 3000 );
config.host = (process.env.IDM_HOST || '180.179.214.135' + config.port);

config.debug = to_boolean(process.env.IDM_DEBUG, true);

// HTTPS enable
config.https = {
    enabled: to_boolean(process.env.IDM_HTTPS_ENABLED, false),
    cert_file: 'certs/idm-2018-cert.pem',
    key_file: 'certs/idm-2018-key.pem',
    ca_certs: [],
    port: (process.env.IDM_HTTPS_PORT || 443 )
};

// Config email list type to use domain filtering
config.email_list_type = (process.env.IDM_EMAIL_LIST || null );   // whitelist or blacklist

// Secret for user sessions in web
config.session = {
    secret:  (process.env.IDM_SESSION_SECRET || require('crypto').randomBytes(20).toString('hex')),       // Must be changed
    expires: (process.env.IDM_SESSION_DURATION || 60 * 60 * 1000)     // 1 hour
}

// Key to encrypt user passwords
config.password_encryption = {
    key: (process.env.IDM_ENCRYPTION_KEY || 'nodejs_idm')   // Must be changed
}

// Enable CORS
config.cors = {
    enabled: to_boolean(process.env.IDM_CORS_ENABLED, false),
    options: {
        /* eslint-disable snakecase/snakecase */
        origin: to_array(process.env.IDM_CORS_ORIGIN, '*'),
        methods: to_array(process.env.IDM_CORS_METHODS, ['GET','HEAD','PUT','PATCH','POST','DELETE']),
        allowedHeaders: (process.env.IDM_CORS_ALLOWED_HEADERS || '*'),
        exposedHeaders: (process.env.IDM_CORS_EXPOSED_HEADERS || undefined),
        credentials: (process.env.IDM_CORS_CREDENTIALS || undefined),
        maxAge: (process.env.IDM_CORS_MAS_AGE || undefined),
        preflightContinue: (process.env.IDM_CORS_PREFLIGHT || false),
        optionsSuccessStatus: (process.env.IDM_CORS_OPTIONS_STATUS || 204)
        /* eslint-enable snakecase/snakecase */
    }
}

// Config oauth2 parameters
config.oauth2 = {
    allow_empty_state: (process.env.IDM_OAUTH_EMPTY_STATE || false),                       // allow empty state in request
    authorization_code_lifetime: (process.env.IDM_OAUTH_AUTH_LIFETIME || 5 * 60),        // Five minutes
    access_token_lifetime: (process.env.IDM_OAUTH_ACC_LIFETIME || 60 * 60),              // One hour
    ask_authorization: (process.env.IDM_OAUTH_ASK_AUTH || true),                         // Prompt a message to users to allow the application to read their details
    refresh_token_lifetime: (process.env.IDM_OAUTH_REFR_LIFETIME || 60 * 60 * 24 * 14),  // Two weeks
    unique_url: (process.env.IDM_OAUTH_UNIQUE_URL || false)                              // This parameter allows to verify that an application with the same url
                                                                                         // does not exist when creating or editing it. If there are already applications
                                                                                         // with the same URL, they should be changed manually
}

// Config api parameters
config.api = {
    token_lifetime: (process.env.IDM_API_LIFETIME || 60*60)     // One hour
}

// Configure Policy Decision Point (PDP)
//  - IdM can perform basic policy checks (HTTP verb + path)
//  - AuthZForce can perform basic policy checks as well as advanced
// If authorization level is advanced you can create rules, HTTP verb+resource and XACML advanced. In addition
// you need to have an instance of authzforce deployed to perform advanced authorization request from a Pep Proxy.
// If authorization level is basic, only HTTP verb+resource rules can be created
config.authorization = {
    level: (process.env.IDM_PDP_LEVEL || 'basic'),     // basic|advanced
    authzforce: {
        enabled: to_boolean(process.env.IDM_AUTHZFORCE_ENABLED, false),
        host: (process.env.IDM_AUTHZFORCE_HOST || 'localhost'),
        port: (process.env.IDM_AUTHZFORCE_PORT||  8080),
    }
}

// Enable usage control and configure where is the Policy Translation Point
config.usage_control = {
    enabled: to_boolean(process.env.IDM_USAGE_CONTROL_ENABLED, false),
    ptp: {
        host: (process.env.IDM_PTP_HOST || 'localhost'),
        port: (process.env.IDM_PTP_PORT||  8080),
    } 
}

// Database info
config.database  = {
    host:     (process.env.IDM_DB_HOST || 'localhost'),
    password: (process.env.IDM_DB_PASS || 'idm'),
    username: (process.env.IDM_DB_USER || 'root'),
    database: (process.env.IDM_DB_NAME || 'idm'),
    dialect:  (process.env.IDM_DB_DIALECT || 'mysql'),
    port:     (process.env.IDM_DB_PORT || undefined)
};

// External user authentication
config.external_auth = {
    enabled: (process.env.IDM_EX_AUTH_ENABLED || false ),
    id_prefix: (process.env.IDM_EX_AUTH_ID_PREFIX || 'external_'),
    password_encryption: (process.env.IDM_EX_AUTH_PASSWORD_ENCRYPTION || 'sha1'),    // bcrypt, sha1 and pbkdf2 supported
    password_encryption_key: (process.env.IDM_EX_AUTH_PASSWORD_ENCRYPTION_KEY || undefined),
    password_encryption_opt: {
        digest: (process.env.IDM_EX_AUTH_PASSWORD_ENCRYPTION_OPT_DIGEST || 'sha256'),
        keylen: (process.env.IDM_EX_AUTH_PASSWORD_ENCRYPTION_OPT_KEYLEN || 64),
        iterations: (process.env.IDM_EX_AUTH_PASSWORD_ENCRYPTION_OPT_ITERATIONS || 27500)
    },
    database: {
        host: (process.env.IDM_EX_AUTH_DB_HOST ||'localhost'),
        port: (process.env.IDM_EX_AUTH_PORT || undefined),
        database: (process.env.IDM_EX_AUTH_DB_NAME ||'db_name'),
        username: (process.env.IDM_EX_AUTH_DB_USER || 'db_user'),
        password: (process.env.IDM_EX_AUTH_DB_PASS ||'db_pass'),
        user_table: (process.env.IDM_EX_AUTH_DB_USER_TABLE ||'user_view'),
        dialect: (process.env.IDM_EX_AUTH_DIALECT || 'mysql')
    }
}

// Email configuration
config.mail = {
    transport: (process.env.IDM_EMAIL_TRANSPORT || 'smtp'),
    domain: (process.env.IDM_EMAIL_DOMAIN || ''),
    host: (process.env.IDM_EMAIL_HOST || 'localhost'),
    port: (process.env.IDM_EMAIL_PORT || 25),
    from: (process.env.IDM_EMAIL_ADDRESS || 'noreply@localhost'),
    mailgun_api_key: (process.env.IDM_MAILGUN_API_KEY || '')
}

// Config themes
config.site = {
    title: (process.env.IDM_TITLE || 'Identity Manager'),
    theme: (process.env.IDM_THEME || 'default')
};

// Config eIDAS Authentication
config.eidas = {
    enabled:       to_boolean(process.env.IDM_EIDAS_ENABLED, false),
    gateway_host:  (process.env.IDM_EIDAS_GATEWAY_HOST || 'localhost'),
    node_host:      (process.env.IDM_EIDAS_NODE_HOST || 'https://se-eidas.redsara.es/EidasNode/ServiceProvider'),
    metadata_expiration: (process.env.IDM_EIDAS_METADATA_LIFETIME || 60 * 60 * 24 * 365) // One year
}

// Enables the possibility of adding identity attributes in users' profile
config.identity_attributes = {
    /* eslint-disable snakecase/snakecase */
    enabled: false,
    attributes: [
        {name: 'Vision', key: 'vision', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Color Perception', key: 'color', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Hearing', key: 'hearing', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Vocal Capability', key: 'vocal', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Manipulation Strength', key: 'manipulation', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Reach', key: 'reach', type: 'number', minVal: '0', maxVal: '100'},
        {name: 'Cognition', key: 'cognition', type: 'number', minVal: '0', maxVal: '100'}
    ]
    /* eslint-enable snakecase/snakecase */
}


if (config.session.secret === 'nodejs_idm' || config.password_encryption.key  === 'nodejs_idm'){
    /* eslint-disable no-console */
    console.log('****************');
    console.log('WARNING: The current encryption keys match the defaults found in the plaintext');
    console.log('         template file - please update for a production instance');
    console.log('****************');
    /* eslint-enable no-console */

}


module.exports = config;
