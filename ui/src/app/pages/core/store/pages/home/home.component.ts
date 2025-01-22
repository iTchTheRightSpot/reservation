import { ChangeDetectionStrategy, Component } from '@angular/core';
import { concatMap, delay, from, of, repeat } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-home',
  imports: [],
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
      [style.background-image]="'url(' + images() + ')'"
    ></div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class HomeComponent {
  protected readonly images = toSignal(
    of([
      './salon-1.jpg',
      './salon-3.jpg',
      './salon-2.jpg'
    ]).pipe(
      concatMap((photos) =>
        from(photos).pipe(
          concatMap((photo) => of(photo).pipe(delay(5000))),
          repeat()
        )
      ),
    ),
    { initialValue: './salon-1.jpg' }
  );
}
