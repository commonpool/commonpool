import {MarkedRenderer} from 'ngx-markdown';
import {DomSanitizer} from '@angular/platform-browser';

export class MarkdownRenderer extends MarkedRenderer {

  public constructor(private sanitize: DomSanitizer) {
    super();
  }

  public link(href: string | null, title: string | null, text: string): string {
    let html = super.link(href, title, text);
    html = html.replace('href', 'routerLink');
    return html;
  }

  text(text: string): string {
    return super.text(text);
  }

  html(html: string): string {
    return super.html(html);
  }


}
