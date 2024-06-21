import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useAuth } from "../context/AuthContext"
import { faGoogle } from "@fortawesome/free-brands-svg-icons";
import { faDoorOpen } from "@fortawesome/free-solid-svg-icons";

const LoginButton = () => {
    const {user, login, logout} = useAuth();

    return (
        <>
            {user === null ?
                <button onClick={login} className="rounded-full ring-2 ring-black px-3 hover:bg-black hover:text-white">
                    <FontAwesomeIcon icon={faGoogle} className="pr-2"/>
                    Join
                </button> :
                <button onClick={logout} className="rounded-full ring-2 ring-black px-3 hover:bg-black hover:text-white">
                    <FontAwesomeIcon icon={faDoorOpen} className="pr-2"/>
                    Leave
                </button>
            }
        </>
    )
}

export default LoginButton
