import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-circle-fill',
  template: `
    <svg [attr.width]="width" [attr.height]="height" viewBox="0 0 16 16" class="bi bi-circle-fill" [attr.fill]="fill"
         xmlns="http://www.w3.org/2000/svg">
      <circle cx="8" cy="8" r="8"/>
    </svg>
  `,
  styles: []
})
// tslint:disable-next-line:component-class-suffix
export class CircleFillIcon implements OnInit {

  constructor() {
  }

  @Input()
  width: any = '1em';

  @Input()
  height: any = '1em';

  @Input()
  fill = 'currentcolor';

  ngOnInit(): void {
  }


}
