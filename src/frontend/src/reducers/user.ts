import {userActions} from '../actions/actions.types';
import {UserStore} from "../store/store.types";
import {UserActions} from "../actions/user";

// visible for test
export const initialState: UserStore = {
  isLogin: false,
  categories: [],
  vocabularyLists: [],
  selectedCategory: {name: "", owner: "", columns: ["", ""]}
};

export const user = (state = initialState, action: UserActions): UserStore => {
  switch (action.type) {
    case userActions.LOGIN:
      return {...state, isLogin: true};
    case userActions.LOGOUT:
      return initialState;
    case userActions.STORE_CATEGORIES:
      return {...state, categories: action.payload};
    case userActions.STORE_VOCABULARY_LISTS:
      return {...state, vocabularyLists: action.payload};
    case userActions.SET_SELECTED_CATEGORY:
      return {...state, selectedCategory: action.payload};
    case userActions.REMOVE_CATEGORY:
      let stateCopy = {...state};
      if (state.selectedCategory.id === action.payload) {
        stateCopy = {...stateCopy, selectedCategory: {name: "", owner: "", columns: ["", ""]}};
      }
      return {
        ...stateCopy, vocabularyLists: stateCopy.vocabularyLists
          .filter(vocabularyList => vocabularyList.categoryId !== action.payload),
        categories: stateCopy.categories.filter(category => category.id !== action.payload)
      };
    default: {
      return state;
    }
  }
};
