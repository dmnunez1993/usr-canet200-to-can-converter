import Errors from "../models/errors";
import { HTTP_NO_CONTENT } from "../utils/constants/http";
import { ERROR, FAILED, SUCCESS } from "../utils/constants/tags";
import { getAbsoluteUrl } from "./urls";

interface ListRequestStatus<Type> {
  status: string;
  detail?: string;
  data?: {
    items: Type[];
    count: number;
  };
  errors?: Errors;
}

export const getList = async <Type>(
  url: string,
  limit: number,
  offset: number,
  additionalParams: Map<string, string | undefined> | undefined = undefined
) => {
  const params = new URLSearchParams();
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());

  if (additionalParams !== undefined) {
    for (let entry of Array.from(additionalParams.entries())) {
      if (entry[1] !== undefined) {
        params.set(entry[0], entry[1]);
      }
    }
  }

  const fullUrl = getAbsoluteUrl(url) + "?" + params.toString();

  let response: Response = { ok: false } as Response;
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "GET",
      credentials: "include",
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }
  const responseData: ListRequestStatus<Type> = await response.json();

  return responseData;
};

export interface ItemRequestStatus<Type> {
  status: string;
  detail?: string;
  data?: Type;
  errors?: Errors;
}

export const createItem = async <Type>(
  url: string,
  item: Type
): Promise<ItemRequestStatus<Type>> => {
  let response: Response = { ok: false } as Response;
  const fullUrl = getAbsoluteUrl(url);
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "POST",
      credentials: "include",
      body: JSON.stringify(item),
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }

  const responseData: ItemRequestStatus<Type> = await response.json();

  return responseData;
};

export const getItem = async <Type>(
  url: string,
  additionalParams: Map<string, string | undefined> | undefined = undefined
): Promise<ItemRequestStatus<Type>> => {
  const params = new URLSearchParams();

  if (additionalParams !== undefined) {
    for (let entry of Array.from(additionalParams.entries())) {
      if (entry[1] !== undefined) {
        params.set(entry[0], entry[1]);
      }
    }
  }

  const fullUrl = getAbsoluteUrl(url) + "?" + params.toString();

  let response: Response = { ok: false } as Response;
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "GET",
      credentials: "include",
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }

  const responseData: ItemRequestStatus<Type> = await response.json();

  return responseData;
};

export const updateItem = async <Type>(
  url: string,
  item: Type
): Promise<ItemRequestStatus<Type>> => {
  const fullUrl = getAbsoluteUrl(url);
  let response: Response = { ok: false } as Response;
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "PUT",
      credentials: "include",
      body: JSON.stringify(item),
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }

  const responseData: ItemRequestStatus<Type> = await response.json();

  return responseData;
};

export const deleteItem = async <Type>(
  url: string,
  item: Type
): Promise<ItemRequestStatus<Type>> => {
  const fullUrl = getAbsoluteUrl(url);
  let response: Response = { ok: false } as Response;
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "DELETE",
      credentials: "include",
      body: JSON.stringify(item),
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }

  if (response.status === HTTP_NO_CONTENT) {
    return { status: SUCCESS };
  }

  return {
    status: ERROR,
    detail: "Hubo un error eliminando el ítem. Intente de vuelta",
  };
};

export const postItem = async <Type>(
  url: string,
  item: Type
): Promise<ItemRequestStatus<Type>> => {
  let response: Response = { ok: false } as Response;
  const fullUrl = getAbsoluteUrl(url);
  try {
    response = await fetch(fullUrl, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "POST",
      credentials: "include",
      body: JSON.stringify(item),
    });
  } catch (error) {
    return {
      status: FAILED,
      detail: "Server not Found",
    };
  }

  const responseData: ItemRequestStatus<Type> = await response.json();

  return responseData;
};
