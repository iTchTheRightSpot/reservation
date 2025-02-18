import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Image } from 'primeng/image';
import { Button } from 'primeng/button';
import { Carousel } from 'primeng/carousel';
import { LightDarkModeService } from '@shared/data-access/light-dark-mode.service';
import { NgOptimizedImage } from '@angular/common';
import { Drawer } from 'primeng/drawer';
import { CreateBookingComponent } from '@shared/ui/reservation/create-booking.component';
import { BookingService } from '@shared/data-access/booking.service';
import { ServiceTypeImpl } from '@shared/data-access/service-type.service';
import { Router } from '@angular/router';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';
import { ABOUT_ROUTE } from '@store/store.routes';
import { ReservationUtilComponent } from '@shared/ui/reservation/reservation-util.component';

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
export class HomeComponent extends ReservationUtilComponent {
  constructor(serviceTypes: ServiceTypeImpl, bookingService: BookingService) {
    super(serviceTypes, bookingService);
  }

  private readonly router = inject(Router);
  protected readonly service = inject(LightDarkModeService);

  protected toggleNewBookings = false;
  protected readonly name = 'Revive Hair Studio';
  protected readonly logo = './logo.jpeg';
  protected readonly img = './salon-1.jpg';
  protected readonly instagram =
    'https://www.instagram.com/revivehairstudiolekki/';
  protected readonly address =
    'Ologolo Lekki, Opposite White Oak Estate, Eti-Osa, Lagos, NG';
  protected readonly imgs = ['./salon-2.jpg', './salon-3.jpg'];
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
}
