export interface StaffScheduleEmitter {
  staff_id: string;
  date: Date;
  page: number;
  size: number;
}

export interface DeleteScheduleModel extends StaffScheduleEmitter {
  schedule_id: number;
}
