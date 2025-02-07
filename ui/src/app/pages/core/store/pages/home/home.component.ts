import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Image } from 'primeng/image';
import { Button } from 'primeng/button';
import { Carousel } from 'primeng/carousel';
import { LightDarkModeService } from '@shared/data-access/light-dark-mode.service';
import { NgOptimizedImage } from '@angular/common';
import { Drawer } from 'primeng/drawer';
import { CreateBookingComponent } from '@shared/ui/create-booking/create-booking.component';
import { EMPTY, map, Subject, switchMap } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import {
  ConfirmModel,
  DateModel,
  StaffModel
} from '@shared/model/shared.model';
import { BookingService } from '@shared/data-access/booking.service';
import { ServiceTypeService } from '@store/pages/reservation/pages/service-type/service-type.service';
import { Router } from '@angular/router';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';
import { ABOUT_ROUTE } from '@store/store.routes';

@Component({
  selector: 'app-home',
  imports: [
    Image,
    Button,
    Carousel,
    NgOptimizedImage,
    Drawer,
    CreateBookingComponent
  ],
  styles: [
    `
      .karla-font {
        font-family: 'Karla', serif;
        font-optical-sizing: auto;
        font-weight: 300;
        font-style: normal;
      }
    `
  ],
  templateUrl: 'home.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class HomeComponent {
  private readonly router = inject(Router);
  protected readonly service = inject(LightDarkModeService);
  private readonly serviceTypes = inject(ServiceTypeService);
  private readonly bookingService = inject(BookingService);

  protected toggleNewBookings = false;
  protected readonly name = 'Revive Hair Studio';
  protected readonly logo = './logo.jpeg';
  protected readonly img = './salon-1.jpg';
  protected readonly instagram = 'https://instagram.com';
  protected readonly address =
    '456 Count Dracula Way, Transylvania Hills, Romania';
  protected readonly imgs = [
    './salon-1.jpg',
    './salon-2.jpg',
    './salon-3.jpg',
    './salon-2.jpg'
  ];
  protected readonly responsiveOptions = [
    {
      breakpoint: '1400px',
      numVisible: 2,
      numScroll: 1
    },
    {
      breakpoint: '1199px',
      numVisible: 3,
      numScroll: 1
    },
    {
      breakpoint: '767px',
      numVisible: 2,
      numScroll: 1
    },
    {
      breakpoint: '575px',
      numVisible: 1,
      numScroll: 1
    }
  ];

  protected readonly about = async () =>
    await this.router.navigate([`${CORE_ROUTE}/${STORE_ROUTE}/${ABOUT_ROUTE}`]);

  protected readonly servicesEmitter = new Subject<boolean>();
  protected readonly services = toSignal<
    ApiResponse<string[]>,
    ApiResponse<string[]>
  >(
    this.servicesEmitter.asObservable().pipe(
      switchMap(b =>
        !b
          ? EMPTY
          : this.serviceTypes.services().pipe(
              map(
                o =>
                  <ApiResponse<string[]>>{
                    state: o.state,
                    message: o.message,
                    data: o.data?.map(s => s.name)
                  }
              )
            )
      )
    ),
    { initialValue: { state: ApiState.LOADED } }
  );

  protected readonly staffEmitter = new Subject<string[]>();
  protected readonly staffs = toSignal(
    this.staffEmitter
      .asObservable()
      .pipe(switchMap(a => this.bookingService.staffsByServiceTypes(a))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<StaffModel[]> }
  );

  protected readonly validDatesEmitter = new Subject<{
    date: Date;
    services: string[];
    staff_id: string;
  }>();
  protected readonly validDates = toSignal(
    this.validDatesEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.validDates(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<DateModel[]> }
  );

  protected readonly reserveBookingEmitter = new Subject<ConfirmModel>();
  protected readonly reserveBooking = toSignal(
    this.reserveBookingEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.create(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}
