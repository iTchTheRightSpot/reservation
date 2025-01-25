export enum BookingStatus {
  CONFIRMED = 'CONFIRMED',
  CANCELLED = 'CANCELLED'
}

export interface BookingsModel {
  reservation_id: number;
  name: string;
  email: string;
  description: string | null;
  phone: string | null;
  price: number;
  status: BookingStatus;
  created_at: number;
  scheduled_for: number;
  expire_at: number;
  services: string[];
}

export interface BookingsRequestPayload {
  user_id: string;
  page: number;
  size: number;
}

export interface UpdateBookingStatusPayload {
  reservation_id: number;
  status: BookingStatus;
}

export const DummyBookingsModels = (num: number) =>
  Array.from({ length: num }, (_, index) => {
    const from = new Date();
    from.setDate(index);
    const to = new Date(from);
    to.setHours(8 * 60 * 60);

    const obj: BookingsModel = {
      reservation_id: index + 1,
      name: `Firstname-${index + 1}`,
      email: `email-${index + 1}@email.com`,
      description:
        index % 2 === 0
          ? 'Lorem ipsum dolor sit amet, consectetur adipisicing elit.'
          : 'null',
      phone: index % 2 === 0 ? '1234567890' : null,
      price: Math.random() * 100,
      status:
        index % 2 === 0 ? BookingStatus.CONFIRMED : BookingStatus.CANCELLED,
      created_at: from.getTime(),
      scheduled_for: from.getTime(),
      expire_at: to.getTime(),
      services: ['erp', 'mens haircut']
    };
    return obj;
  });
