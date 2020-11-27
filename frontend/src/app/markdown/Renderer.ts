import {MarkedOptions, MarkedRenderer} from 'ngx-markdown';
import {Sanitizer, SecurityContext} from '@angular/core';
import {DomSanitizer} from '@angular/platform-browser';
//
// export class MarkdownRenderer extends MarkedRenderer {
//
//   public constructor(private sanitize: DomSanitizer) {
//     super();
//   }
//
//   public link(href: string | null, title: string | null, text: string): string {
//     let html = super.link(href, title, text);
//     html = html.replace('href', 'routerLink');
//     return html;
//   }
//
//   text(text: string): string {
//     console.log(text);
//
//     text = text.replace(/:::user:::\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b:::/, '<a href="#">USER</a>');
//
//     return super.text(text);
//   }
//
// }
