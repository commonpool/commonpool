import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
import {UserProfileComponent} from './user-profile/user-profile.component';
import {ConversationThreadListComponent} from './chat/conversation-thread-list/conversation-thread-list.component';
import {ResourceInquiryComponent} from './resources/resource-inquiry/resource-inquiry.component';
import {ConversationThreadComponent} from './chat/conversation-thread/conversation-thread.component';
import {OfferListComponent} from './offers/offer-list/offer-list.component';
import {OfferDetailsComponent} from './offers/offer-details/offer-details.component';
import {CreateOfferComponent} from './offers/create-offer/create-offer.component';


const routes: Routes = [
  {
    path: '',
    component: ResourceListViewComponent
  }, {
    path: 'resources/new',
    component: CreateOrEditResourceComponent
  }, {
    path: 'resources/:id',
    component: ResourceDetailsComponent
  }, {
    path: 'resources/:id/inquire',
    component: ResourceInquiryComponent
  }, {
    path: 'resources/:id/edit',
    component: CreateOrEditResourceComponent
  }, {
    path: 'profiles/:id',
    component: UserProfileComponent
  }, {
    path: 'messages',
    component: ConversationThreadListComponent
  }, {
    path: 'messages/:id',
    component: ConversationThreadComponent
  }, {
    path: 'offers',
    component: OfferListComponent
  }, {
    path: 'offers/new',
    component: CreateOfferComponent,
  }, {
    path: 'offers/:id',
    component: OfferDetailsComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
