export const emptyValueOnUndefined = <Type>(value?: Type) => {
  if (value === undefined) {
    return "";
  }

  return value;
};

export const emptyValueOnNull = <Type>(value: Type | null) => {
  if (value === null) {
    return "";
  }

  return value;
};

export const undefinedOnNull = <Type>(value: Type | null) => {
  if (value === null) {
    return undefined;
  }

  return value;
};
