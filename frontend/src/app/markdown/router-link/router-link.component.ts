import {ChangeDetectionStrategy, Component, Input, OnInit} from '@angular/core';
import {DomSanitizer, SafeHtml} from '@angular/platform-browser';

// @Component({
//   template: '<a [routerLink]="href">{{text}}</a>',
//   changeDetection: ChangeDetectionStrategy.OnPush
// })
// export class RouterLinkComponent {
//   constructor(private sanitizer: DomSanitizer) {
//
//   }
//
//   private _href: SafeHtml;
//   @Input()
//   set href(value: SafeHtml | string) {
//     this._href = value
//   }
//
//   get href(): SafeHtml {
//     return this._href;
//   }
//
//   @Input() public text: string;
// }
