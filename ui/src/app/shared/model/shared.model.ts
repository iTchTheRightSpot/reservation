export interface ConfirmModel {
  staff_id: string;
  name: string;
  email: string;
  description: string;
  phone: string;
  services: string[];
  timezone: string;
  time: string;
}

export interface DateModel {
  date: string;
  times: string[];
}

export function DummyDates(month: number, year: number) {
  const results: DateModel[] = [];
  const startOfMonth = new Date(year, month, 1);
  const endOfMonth = new Date(year, month, 0);

  for (let day = startOfMonth.getDate(); day <= endOfMonth.getDate(); day++) {
    if (day % 2 === 0) {
      const times: string[] = [];
      for (let hour = 9; hour <= 17; hour++)
        times.push(new Date(year, month, day, hour).getTime().toString());

      results.push({
        date: new Date(year, month, day).getTime().toString(),
        times
      });
    }
  }
  return results;
}

export interface StaffModel {
  staff_id: string;
  name: string;
  image_key: string | null;
  bio: string;
}

export const DummyStaffModels: StaffModel[] = [
  {
    staff_id: '1',
    name: 'Tony',
    image_key: './salon-1.jpg',
    bio: 'Ready to put a smile on your face🌞'
  },
  {
    staff_id: '2',
    name: 'benjamin',
    image_key: './salon-2.jpg',
    bio: 'Ready to put a smile on your face🌞'
  },
  {
    staff_id: '3',
    name: 'phil',
    image_key: './salon-3.jpg',
    bio: 'Ready to put a smile on your face🌞'
  }
];
