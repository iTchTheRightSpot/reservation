import * as moment from 'moment-timezone';
import { TIMEZONE } from '@root/app.util';

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

export const filterValidDatesFromDatesInAMonth = (
  date: Date,
  valid: DateModel[]
) => {
  const validDates = valid.map(d =>
    moment.tz(Number(d.date), TIMEZONE).toDate()
  );

  const endOfMonth = new Date(date.getFullYear(), date.getMonth() + 1, 0);
  const daysInMonth = endOfMonth.getDate();

  const allDatesInMonth = Array.from(
    { length: daysInMonth },
    (_, index) => new Date(date.getFullYear(), date.getMonth(), index + 1)
  );

  return allDatesInMonth.filter(
    date =>
      !validDates.some(validDate => {
        return (
          validDate.getDate() === date.getDate() &&
          validDate.getMonth() === date.getMonth() &&
          validDate.getFullYear() === date.getFullYear()
        );
      })
  );
};

export const findDatesInDateModel = (d: Date, arr: DateModel[]) => {
  const find = arr.find(obj => {
    const t = moment.tz(Number(obj.date), TIMEZONE).toDate();
    return (
      t.getDate() === d.getDate() &&
      t.getMonth() === d.getMonth() &&
      t.getFullYear() === d.getFullYear()
    );
  });
  if (!find) return undefined;
  return find.times.map(a => ({
    original: a,
    format: moment.tz(Number(a), TIMEZONE).format('h:mm a')
  }));
};

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

export interface ServiceTypeModel {
  name: string;
  price: number;
  duration: number;
}

export const DummyServiceTypes: ServiceTypeModel[] = [
  {
    name: 'power grooming',
    price: 35,
    duration: 100
  },
  {
    name: 'overgrown lawns',
    price: 30,
    duration: 50
  },
  {
    name: 'utility cuts',
    price: 840,
    duration: 7200
  },
  {
    name: 'weekly trim and mow',
    price: 100,
    duration: 86400
  },
  {
    name: 'pre-call service',
    price: 30,
    duration: 3600
  }
];
