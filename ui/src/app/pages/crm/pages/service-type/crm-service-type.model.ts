export interface CRM_ServiceTypeModel {
  name: string;
  price: number;
  is_visible: boolean;
  is_reoccurring: boolean;
  duration: number;
  clean_up_time: number;
}

export const CRM_DummyServiceTypes = (num: number): CRM_ServiceTypeModel[] =>
  Array.from(
    { length: num },
    (_, index) =>
      ({
        name: `Service Type ${index + 1}`,
        price: Math.floor(Math.random() * 100) + 1,
        is_visible: Math.random() > 0.5,
        is_reoccurring: Math.random() > 0.5,
        duration: Math.floor(Math.random() * 1000) + 1,
        clean_up_time: Math.floor(Math.random() * 500) + 1
      }) as CRM_ServiceTypeModel
  );
