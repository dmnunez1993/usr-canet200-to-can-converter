export interface UsrCanetDeviceConverterConfig {
  host?: string;
  port?: number;
  target?: string;
}

export interface UsrCanetConverterConfig {
  deviceConverters: UsrCanetDeviceConverterConfig[];
}

export const newUsrCanetConverterConfig = () => {
  return {
    deviceConverters: [],
  };
};
