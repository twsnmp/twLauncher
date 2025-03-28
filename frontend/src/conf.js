export const makeTwsnmpfcParams = (conf) => {
  const params = [];
  params.push("-datastore");
  params.push(conf.DataStore);
  if (conf.Password) {
    params.push("-password");
    params.push(conf.Password);
  }
  params.push("-port");
  params.push(conf.Port + "");
  if (conf.Local) {
    params.push("-local");
  } else if (conf.TLS) {
    params.push("-tls");
  }
  return params;
};

export const makeTwpcapParams = (conf) => {
  const params = [];
  params.push("-syslog");
  params.push(conf.Syslog);
  params.push("-iface");
  params.push(conf.Iface);
  params.push("-interval");
  params.push(conf.Interval + "");
  params.push("-retention");
  params.push(conf.Retention + "");
  return params;
};

export const makeTwWifiScanParams = (conf) => {
  const params = [];
  params.push("-syslog");
  params.push(conf.Syslog);
  params.push("-interval");
  params.push(conf.Interval + "");
  return params;
};

export const makeTwWinLogParams = (conf) => {
  const params = [];
  params.push("-syslog");
  params.push(conf.Syslog);
  params.push("-interval");
  params.push(conf.Interval + "");
  if( conf.Remote) {
    params.push("-remote");
    params.push(conf.Remote);
    params.push("-user");
    params.push(conf.User);
    params.push("-password");
    params.push(conf.Passowd);
  }
  return params;
};

export const makeTwLogEyeParams = (conf) => {
  const params = [];
  params.push("start");
  params.push("--apiPort");
  params.push(conf.APIPort + "");
  params.push("--config");
  params.push(conf.ConfFile);
  if (conf.Key && conf.Cert) {
    params.push("--serverKey");
    params.push(conf.Key);
    params.push("--serverCert");
    params.push(conf.Cert);
  }
  if (conf.CACert) {
    params.push("--caCert");
    params.push(conf.CACert);
  }
  console.log(params);
  return params;
};


export const getConfFromParams = (info,params, task) => {
  const conf = {
    DataStore: "",
    Password: "",
    Local: false,
    Port: 8080,
    APIPort: 8081,
    ConfFile: "",
    Key: "",
    Cert: "",
    CACert: "",
    Syslog: "",
    Interval: 600,
    Retention: 3600,
    Iface: "",
    Remote: "",
    User: "",
    Task: task && info.Env == "windows",
  };
  for (let i = 0; i < params.length; i++) {
    switch (params[i]) {
      case "-datastore":
        if (i < params.length - 1) {
          conf.DataStore = params[i + 1];
          i++;
        }
        break;
      case "-password":
        if (i < params.length - 1) {
          conf.Password = params[i + 1];
          i++;
        }
        break;
      case "-port":
        if (i < params.length - 1) {
          conf.Port = params[i + 1] * 1;
          i++;
        }
        break;
      case "-local":
        conf.Local = true;
        break;
      case "-syslog":
        if (i < params.length - 1) {
          conf.Syslog = params[i + 1];
          i++;
        }
        break;
      case "-iface":
        if (i < params.length - 1) {
          conf.Iface = params[i + 1];
          i++;
        }
        break;
      case "-interval":
        if (i < params.length - 1) {
          conf.Interval = params[i + 1] * 1;
          i++;
        }
        break;
      case "-retention":
        if (i < params.length - 1) {
          conf.Retention = params[i + 1] * 1;
          i++;
        }
        break;
      case "-remote":
        if (i < params.length - 1) {
          conf.Remote = params[i + 1];
          i++;
        }
        break;
      case "-user":
        if (i < params.length - 1) {
          conf.User = params[i + 1];
          i++;
        }
      case "--apiPort":
        if (i < params.length - 1) {
          conf.APIPort = params[i + 1] * 1;
          i++;
        }
        break;
      case "--config":
        if (i < params.length - 1) {
          conf.ConfFile = params[i + 1];
          i++;
        }
        break;
      case "--serverKey":
        if (i < params.length - 1) {
          conf.Key = params[i + 1];
          i++;
        }
        break;
      case "--serverCert":
        if (i < params.length - 1) {
          conf.Cert = params[i + 1];
          i++;
        }
        break;
      case "--caCert":
        if (i < params.length - 1) {
          conf.CACert = params[i + 1];
          i++;
        }
        break;
    }
  }
  if (conf.Interval < 60) {
    conf.Interval = 600;
  }
  if (conf.Retention < 600) {
    conf.Retention = 3600;
  }
  return conf;
};
