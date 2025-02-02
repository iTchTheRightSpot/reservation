export interface Schedule {
  schedule_id: number;
  start: string;
  end: string;
  is_visible: boolean;
  is_reoccurring: boolean;
}

export interface CreateScheduleModel {
  staff_id: string;
  start: Date;
  duration: number;
  is_visible: boolean;
  is_reoccurring: boolean;
}

export interface UpdateScheduleModel {
  staff_id: string;
  schedule_id: number;
  is_visible: boolean;
  is_reoccurring: boolean;
}

export const DummySchedules = (num: number): Schedule[] =>
  Array.from({ length: num }, (_, index) => {
    const from = new Date();
    from.setDate(index);
    const to = new Date(from);
    to.setHours(8 * 60 * 60);
    return {
      schedule_id: index + 1,
      start: from.getTime() + '',
      end: to.getTime() + '',
      is_visible: Math.random() > 0.5,
      is_reoccurring: Math.random() > 0.5
    } as Schedule;
  });
