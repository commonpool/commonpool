import {
  ApplicationRef,
  ComponentFactoryResolver,
  Directive,
  ElementRef,
  HostListener,
  Inject,
  Injector
} from '@angular/core';
import {DOCUMENT} from '@angular/common';
import {RouterLinkComponent} from './router-link/router-link.component';
import {DomSanitizer} from '@angular/platform-browser';

// @Directive({
//   // tslint:disable-next-line:directive-selector
//   selector: 'markdown,[markdown]'
// })
// export class ConvertLinksDirective {
//
//   constructor(
//     @Inject(DOCUMENT) private document: Document,
//     private injector: Injector,
//     private applicationRef: ApplicationRef,
//     private componentFactoryResolver: ComponentFactoryResolver,
//     private element: ElementRef<HTMLElement>,
//   ) {
//   }
//
//   @HostListener('ready')
//   public processAnchors() {
//     this.element.nativeElement.querySelectorAll(
//       'a[routerLink]'
//     ).forEach(a => {
//       const parent = a.parentElement;
//       if (parent) {
//         const container = this.document.createElement('span');
//         const component = this.componentFactoryResolver.resolveComponentFactory(
//           RouterLinkComponent
//         ).create(this.injector, [], container);
//         this.applicationRef.attachView(component.hostView);
//         component.instance.href = a.getAttribute('routerLink') || '';
//         component.instance.text = a.textContent || '';
//         parent.replaceChild(container, a);
//       }
//     });
//   }
//
// }
