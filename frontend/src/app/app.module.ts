import {BrowserModule, DomSanitizer} from '@angular/platform-browser';
import {NgModule, SecurityContext} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import {TopNavComponent} from './top-nav/top-nav.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {AuthService} from './auth.service';
import {AppHttpInterceptor} from './api/backend.service';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
import {ConversationThreadComponent} from './chat/conversation-thread/conversation-thread.component';
import {ConversationThreadListComponent} from './chat/conversation-thread-list/conversation-thread-list.component';
import {ScrollingModule} from '@angular/cdk/scrolling';
import {ResourceInquiryComponent} from './resources/resource-inquiry/resource-inquiry.component';
import {OfferListComponent} from './offers/offer-list/offer-list.component';
import {CreateOfferComponent} from './offers/create-offer/create-offer.component';
import {OfferDetailsComponent} from './offers/offer-details/offer-details.component';
import {UserPickerComponent} from './shared/user-picker/user-picker.component';
import {ResourcePickerComponent} from './shared/resource-picker/resource-picker.component';
import {NgSelectModule} from '@ng-select/ng-select';
import {PlusSquareIcon} from './icons/plus-square/plus-square.icon';
import {PlusIcon} from './icons/plus/plus.icon';
import {UsernameComponent} from './shared/username/username.component';
import {ResourceNameComponent} from './shared/resource-name/resource-name.component';
import {TrashIcon} from './icons/trash/trash.icon';
import {CreateOrEditGroupComponent} from './groups/create-or-edit-group/create-or-edit-group.component';
import {RequiredIndicatorComponent} from './shared/required-indicator/required-indicator.component';
import {GroupViewComponent} from './groups/group-view/group-view.component';
import {BoxSeamIcon} from './icons/box-seam/box-seam.icon';
import {PersonIcon} from './icons/person/person.icon';
import {PeopleIcon} from './icons/people/people.icon';
import {PentagonIcon} from './icons/pentagon/pentagon.icon';
import {AsterixIcon} from './icons/asterix/asterix.icon';
import {GroupResourcesViewComponent} from './groups/group-resources-view/group-resources-view.component';
import {GroupMembersViewComponent} from './groups/group-members-view/group-members-view.component';
import {UserOffersViewComponent} from './users/user-offers-view/user-offers-view.component';
import {UserResourcesViewComponent} from './users/user-resources-view/user-resources-view.component';
import {UserGroupsViewComponent} from './users/user-groups-view/user-groups-view.component';
import {UserViewComponent} from './users/user-view/user-view.component';
import {MailboxIcon} from './icons/mailbox/mailbox.icon';
import {EnvelopeIcon} from './icons/envelope/envelope.icon';
import {CheckIcon} from './icons/check/check.icon';
import {CrossIcon} from './icons/cross/cross.icon';
import {DoorOpenIcon} from './icons/door-open/door-open.icon';
import {ArrowRightIcon} from './icons/arrow-right/arrow-right.icon';
import {ArrowLeftIcon} from './icons/arrow-left/arrow-left.icon';
import {GroupInvitesViewComponent} from './groups/group-invites-view/group-invites-view.component';
import {IncomingInvitationComponent} from './users/incoming-invitation/incoming-invitation.component';
import {KickOrLeaveGroupButtonComponent} from './groups/kick-or-leave-group-button/kick-or-leave-group-button.component';
import {CircleFillIcon} from './icons/circle-fill/circle-fill.icon';
import {GroupLinkComponent} from './shared/group-link/group-link.component';
import {ResourceLinkComponent} from './shared/resource-link/resource-link.component';
import {AcceptInvitationComponent} from './groups/accept-invitation/accept-invitation.component';
import {HomePageComponent} from './home/home-page/home-page.component';
import {MarkdownModule, MarkedOptions} from 'ngx-markdown';
import {BlocksComponent} from './chat/blocks/blocks.component';
import {TextObjectComponent} from './chat/text-object/text-object.component';
import {ButtonComponent} from './chat/button/button.component';
import {BlockComponent} from './chat/block/block.component';
import {AttachmentComponent} from './chat/attachment/attachment.component';
import {SectionBlockComponent} from './chat/section-block/section-block.component';
import {ContextBlockComponent} from './chat/context-block/context-block.component';
import {ActionsBlockComponent} from './chat/actions-block/actions-block.component';
import {DividerBlockComponent} from './chat/divider-block/divider-block.component';
import {ImageBlockComponent} from './chat/image-block/image-block.component';
import {HeaderBlockComponent} from './chat/header-block/header-block.component';
import {MessageGroupComponent} from './chat/message-group/message-group.component';
import {SampleComponent} from './sample/sample/sample.component';
import {MarkdownRenderer} from './markdown/renderer';
import {ConvertLinksDirective} from './markdown/convert-links.directive';
import {RouterLinkComponent} from './markdown/router-link/router-link.component';
import {UserLinkComponent} from './shared/user-link/user-link.component';
import {ResourceLink2Component} from './shared/resource-link2/resource-link2.component';
import {TradingHistoryComponent} from './trading/history/trading-history.component';
import {CommonModule} from '@angular/common';
import {GroupOrUserPickerComponent} from './shared/group-or-user-picker/group-or-user-picker.component';
import {MembershipPickerComponent} from './shared/membership-picker/membership-picker.component';
import {OfferListItemComponent} from './offers/offer-list-item/offer-list-item.component';
import {TargetComponent} from './shared/target/target.component';
import {OfferActionComponent} from './offers/offer-action/offer-action.component';
import {DurationComponent} from './shared/duration/duration.component';
import {GroupComponent} from './shared/group/group.component';
import {ValuesComponent} from './scratch/values.component';
import {ValueComponent} from './scratch/value.component';
import {ValueThresholdComponent} from './scratch/value-threshold.component';
import {DimensionValueComponent} from './scratch/dimension-value.component';
import {ResourceEvaluationComponent} from './resources/resource-evaluation/resource-evaluation.component';
import {GroupHistoryComponent} from './groups/history/group-history.component';

const rendererFactory = (domSanitizer: DomSanitizer): MarkedOptions => {
  return {
    renderer: new MarkdownRenderer(domSanitizer),
  };
};

@NgModule({
  declarations: [
    AppComponent,
    ResourceListViewComponent,
    TopNavComponent,
    CreateOrEditResourceComponent,
    ResourceDetailsComponent,
    ConversationThreadComponent,
    ConversationThreadListComponent,
    ResourceInquiryComponent,
    OfferListComponent,
    CreateOfferComponent,
    OfferDetailsComponent,
    UserPickerComponent,
    ResourcePickerComponent,
    PlusSquareIcon,
    PlusIcon,
    UsernameComponent,
    ResourceNameComponent,
    TrashIcon,
    CreateOrEditGroupComponent,
    RequiredIndicatorComponent,
    GroupViewComponent,
    BoxSeamIcon,
    PersonIcon,
    PeopleIcon,
    PentagonIcon,
    AsterixIcon,
    GroupResourcesViewComponent,
    GroupMembersViewComponent,
    UserOffersViewComponent,
    UserResourcesViewComponent,
    UserGroupsViewComponent,
    UserViewComponent,
    MailboxIcon,
    EnvelopeIcon,
    CheckIcon,
    CrossIcon,
    DoorOpenIcon,
    ArrowRightIcon,
    ArrowLeftIcon,
    GroupInvitesViewComponent,
    IncomingInvitationComponent,
    KickOrLeaveGroupButtonComponent,
    CircleFillIcon,
    GroupLinkComponent,
    ResourceLinkComponent,
    AcceptInvitationComponent,
    HomePageComponent,
    BlocksComponent,
    TextObjectComponent,
    ButtonComponent,
    BlockComponent,
    AttachmentComponent,
    SectionBlockComponent,
    ContextBlockComponent,
    ActionsBlockComponent,
    DividerBlockComponent,
    ImageBlockComponent,
    HeaderBlockComponent,
    MessageGroupComponent,
    SampleComponent,
    ConvertLinksDirective,
    RouterLinkComponent,
    UserLinkComponent,
    ResourceLink2Component,
    TradingHistoryComponent,
    GroupOrUserPickerComponent,
    MembershipPickerComponent,
    OfferListItemComponent,
    TargetComponent,
    OfferActionComponent,
    DurationComponent,
    GroupComponent,
    ValuesComponent,
    ValueComponent,
    ValueThresholdComponent,
    DimensionValueComponent,
    ResourceEvaluationComponent,
    GroupHistoryComponent
  ],
  imports: [
    CommonModule,
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    ScrollingModule,
    NgSelectModule,
//    MarkdownModule.forRoot()
    MarkdownModule.forRoot({
      sanitize: SecurityContext.NONE,
      markedOptions: {
        provide: MarkedOptions,
        useFactory: rendererFactory,
        deps: [DomSanitizer]
      }
    })
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: AppHttpInterceptor, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(authService: AuthService) {
    authService.checkLoggedIn();
  }
}

