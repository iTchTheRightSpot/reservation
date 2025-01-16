import { ChangeDetectionStrategy, Component } from '@angular/core';
import { concatMap, delay, from, of, repeat, startWith } from 'rxjs';
import { AsyncPipe } from '@angular/common';

@Component({
  selector: 'app-home',
  imports: [AsyncPipe],
  styles: [
    `
      .trans {
        transition: all 2s ease;
        transition-delay: 1s, 250ms;
      }
    `
  ],
  template: `
    <div
      class="trans h-screen bg-center bg-no-repeat bg-cover"
      [style.background-image]="'url(' + (image$ | async) + ')'"
    ></div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class HomeComponent {
  protected readonly image$ = of([
    './salon-1.jpg',
    './salon-3.jpg',
    './salon-2.jpg'
  ]).pipe(
    concatMap((photos: string[]) =>
      from(photos).pipe(
        concatMap((photo: string) => of(photo).pipe(delay(5000))),
        repeat()
      )
    ),
    startWith('./salon-1.jpg')
  );
}
