interface Errors {
  [key: string]: string[] | Errors | Errors[];
}

export const getFieldErrors = (
  fieldName: string,
  errors?: Errors
): string[] | Errors | Errors[] => {
  if (errors) {
    if (fieldName in errors) {
      return errors[fieldName];
    }
  }

  return [];
};

export const extractArrayErrors = (
  accessor: string,
  errors?: Errors
): Errors[] => {
  if (errors === undefined) {
    return [];
  }

  let extractedErrors: Errors[] = [];

  const hasCorrectFormat = (errorKey: string): boolean => {
    if (!errorKey.startsWith(accessor + "[")) {
      return false;
    }

    if (errorKey.indexOf("]") === -1) {
      return false;
    }

    if (errorKey.indexOf(",", accessor.length) > -1) {
      return false;
    }

    const startIndex = errorKey.indexOf("[");

    if (startIndex === -1) {
      return false;
    }

    const endIndex = errorKey.indexOf("]", startIndex + 1);

    if (endIndex === -1) {
      return false;
    }

    const indexStr = errorKey.substring(startIndex + 1, endIndex);

    if (isNaN(parseInt(indexStr))) {
      return false;
    }

    if (errorKey.indexOf(".", endIndex) === -1) {
      return false;
    }

    return true;
  };

  for (let errorKey in errors) {
    if (!hasCorrectFormat(errorKey)) {
      continue;
    }

    const startIndex = errorKey.indexOf("[");
    const endIndex = errorKey.indexOf("]");
    const itemIndexStr = errorKey.substring(startIndex + 1, endIndex);
    const itemIndex = parseInt(itemIndexStr);

    if (extractedErrors.length < itemIndex + 1) {
      extractedErrors = [];

      for (let idx = 0; idx <= itemIndex; idx++) {
        extractedErrors.push({});
      }
    }
  }

  for (let errorKey in errors) {
    if (!hasCorrectFormat(errorKey)) {
      continue;
    }

    const startIndex = errorKey.indexOf("[");
    const endIndex = errorKey.indexOf("]");
    const itemIndexStr = errorKey.substring(startIndex + 1, endIndex);
    const itemIndex = parseInt(itemIndexStr);

    const dotIndex = errorKey.indexOf(".");

    const attributeName = errorKey.substring(dotIndex + 1, errorKey.length);

    const errorsToAdd = errors[errorKey] as string[];

    const itemErrors = extractedErrors[itemIndex];

    itemErrors[attributeName] = errorsToAdd;
  }

  return extractedErrors;
};

export const errorsAreSame = (a: Errors, b: Errors) => {
  for (let errorKey in a) {
    if (!(errorKey in b)) {
      return false;
    }

    const aErrors = a[errorKey];
    const bErrors = b[errorKey];

    if (!Array.isArray(aErrors) && !Array.isArray(bErrors)) {
      const errorDictsAreTheSame = errorsAreSame(aErrors, bErrors);

      if (!errorDictsAreTheSame) {
        return false;
      }
    } else if (!Array.isArray(aErrors)) {
      return false;
    } else if (!Array.isArray(bErrors)) {
      return false;
    } else {
      const aErrorsString = aErrors as string[];
      const bErrorsString = bErrors as string[];
      for (let error of aErrorsString) {
        const bIndex = bErrorsString.indexOf(error);

        if (bIndex === -1) {
          return false;
        }
      }
    }
  }

  for (let errorKey in b) {
    if (!(errorKey in a)) {
      return false;
    }

    const aErrors = a[errorKey];
    const bErrors = b[errorKey];

    if (!Array.isArray(aErrors) && !Array.isArray(bErrors)) {
      const errorDictsAreTheSame = errorsAreSame(aErrors, bErrors);

      if (!errorDictsAreTheSame) {
        return false;
      }
    } else if (!Array.isArray(aErrors)) {
      return false;
    } else if (!Array.isArray(bErrors)) {
      return false;
    } else {
      const aErrorsString = aErrors as string[];
      const bErrorsString = bErrors as string[];
      for (let error of aErrorsString) {
        const bIndex = bErrorsString.indexOf(error);

        if (bIndex === -1) {
          return false;
        }
      }
    }
  }

  return true;
};

export default Errors;
