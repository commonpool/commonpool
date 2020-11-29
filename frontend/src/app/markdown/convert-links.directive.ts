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
import {UserLinkComponent} from '../shared/user-link/user-link.component';
import {ResourceLink2Component} from '../shared/resource-link2/resource-link2.component';

@Directive({
  // tslint:disable-next-line:directive-selector
  selector: 'markdown,[markdown]'
})
export class ConvertLinksDirective {

  constructor(
    @Inject(DOCUMENT) private document: Document,
    private injector: Injector,
    private applicationRef: ApplicationRef,
    private componentFactoryResolver: ComponentFactoryResolver,
    private element: ElementRef<HTMLElement>,
  ) {

  }

  @HostListener('ready')
  public processAnchors() {

    this.element.nativeElement.querySelectorAll(
      'commonpool-resource'
    ).forEach(a => {
      const parent = a.parentElement;
      if (parent) {
        const container = this.document.createElement('span');
        const component = this.componentFactoryResolver.resolveComponentFactory(
          ResourceLink2Component
        ).create(this.injector, [], container);
        this.applicationRef.attachView(component.hostView);
        component.instance.id = a.getAttribute('id');
        parent.replaceChild(container, a);
      }
    });

    this.element.nativeElement.querySelectorAll(
      'commonpool-user'
    ).forEach(a => {
      const parent = a.parentElement;
      if (parent) {
        const container = this.document.createElement('span');
        const component = this.componentFactoryResolver.resolveComponentFactory(
          UserLinkComponent
        ).create(this.injector, [], container);
        this.applicationRef.attachView(component.hostView);
        component.instance.id = a.getAttribute('id');
        parent.replaceChild(container, a);
      }
    });

  }

}
