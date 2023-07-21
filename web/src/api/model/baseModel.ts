export interface BasicPageParams {
  page?: number;
  page_size?: number;
}

export interface BasicFetchResult<T> {
  items: T[];
  total: number;
}


export interface OptionItem {
  text: string;
  value: string | number;
  status?: number;
  other?: any;
}

export type GetOptionItemsModel = {
  total: number;
  options: OptionItem[]
};
