import React, {lazy, useEffect, useState} from 'react';
import {Route, Switch} from 'react-router-dom';
import {get} from "./utility/restCaller";
import {ToastContainer} from "react-toastify";
import NavigationBar from './components/navigation/navigationBar/NavigationBar';
import Home from './components/Home';
import ProtectedRoute from "./components/navigation/route/ProtectedRoute";
import LoginView from "./containers/login/LoginView";
import {Category} from "./containers/category/CategoryView";
import {useDispatch, useSelector} from "react-redux";
import * as userActions from "./actions/user";
import {AppStore} from "./store/store.types";

const CategoryView = lazy(() => import('./containers/category/CategoryView'));
const VocabularyListView = lazy(() => import('./containers/vocabulary/VocabularyListView'));
const VocabularyView = lazy(() => import('./containers/vocabulary/VocabularyView'));
const ProfileView = lazy(() => import('./containers/profile/ProfileView'));
const TestSettings = lazy(() => import('./containers/learn/test/TestSettings'));
const TestView = lazy(() => import('./containers/learn/test/TestView'));
const App = (): JSX.Element => {
    const isLogin = useSelector((store: AppStore) => store.user.isLogin);
    const dispatch = useDispatch();
    const [timer, setTimer] = useState<NodeJS.Timeout | null>(null);
    useEffect(() => {
        get<{ login: boolean }>('/check-login', (r) => {
            if (r.login) {
                dispatch(userActions.login());
            }
        });
    }, []);
    useEffect(() => {
        if (isLogin) {
            setTimer(setInterval(() => get<{ login: boolean }>('/refresh-token', (r) => {
                if (r.login) {
                    dispatch(userActions.login());
                } else {
                    dispatch(userActions.logout());
                }
            }), 450000));//7,5 minutes
            get<Category[] | null>('/category', data => {
                if (data) {
                    dispatch(userActions.storeCategories(data));
                }
            });
        } else if (timer) {
            clearInterval(timer);
        }
    }, [isLogin]);
    return (
        <>
            <ToastContainer/>
            <NavigationBar/>
            <Switch>
                <ProtectedRoute path='/profile' isAllowed={isLogin}
                                render={(props) => <ProfileView {...props}/>}/>
                <ProtectedRoute path='/vocabulary/:user/:category/:listId' isAllowed={isLogin}
                                render={(props) => <VocabularyView {...props}/>}/>
                <ProtectedRoute path='/vocabulary/:user/:category' isAllowed={isLogin}
                                render={(props) => <VocabularyListView {...props}/>}/>
                <ProtectedRoute path='/vocabulary' isAllowed={true}
                                render={(props) => <CategoryView {...props}/>}/>
                <ProtectedRoute path='/learn/test' isAllowed={true}
                                render={(props) => <TestView {...props}/>}/>
                <ProtectedRoute path='/learn/:user/:category' isAllowed={true}
                                render={(props) => <TestSettings {...props}/>}/>
                <ProtectedRoute path='/learn' isAllowed={true}
                                render={(props) => <CategoryView {...props}/>}/>
                <ProtectedRoute path='/login' isAllowed={!isLogin}
                                render={(props) => <LoginView {...props}/>}/>
                <Route path='/' component={Home} exact/>
            </Switch>
        </>);
};
export default App;
