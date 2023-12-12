import { v4 as uuidv4 } from "uuid";

export interface UniqueListItem<Type> {
  uuid: string;
  item: Type;
}

const newUniqueListItem = <Type>(item: Type): UniqueListItem<Type> => {
  return {
    uuid: uuidv4(),
    item,
  };
};

const listToUniqueListItems = <Type>(items: Type[]): UniqueListItem<Type>[] => {
  const uniqueItems: UniqueListItem<Type>[] = [];

  for (const item of items) {
    uniqueItems.push(newUniqueListItem(item));
  }

  return uniqueItems;
};

const uniqueListToListItems = <Type>(
  uniqueItems: UniqueListItem<Type>[]
): Type[] => {
  const items: Type[] = [];

  for (const uniqueItem of uniqueItems) {
    items.push(uniqueItem.item);
  }

  return items;
};

export { newUniqueListItem, listToUniqueListItems, uniqueListToListItems };
